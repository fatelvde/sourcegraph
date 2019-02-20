import { range } from 'lodash'
import React from 'react'
import VisibilitySensor from 'react-visibility-sensor'
import { combineLatest, Observable, Subject, Subscription } from 'rxjs'
import { catchError, filter, switchMap } from 'rxjs/operators'
import { highlightNode } from '../util/dom'
import { asError, ErrorLike, isErrorLike } from '../util/errors'
import { Repo } from '../util/url'

export interface FetchFileCtx {
    repoName: string
    commitID: string
    filePath: string
    disableTimeout?: boolean
    isLightTheme: boolean
}

interface Props extends Repo {
    commitID: string
    filePath: string
    // How many extra lines to show in the excerpt before/after the ref.
    context?: number
    highlightRanges: HighlightRange[]
    className?: string
    isLightTheme: boolean
    fetchHighlightedFileLines: (ctx: FetchFileCtx, force?: boolean) => Observable<string[]>
}

interface HighlightRange {
    /**
     * The 0-based line number that this highlight appears in
     */
    line: number
    /**
     * The 0-based character offset to start highlighting at
     */
    character: number
    /**
     * The number of characters to highlight
     */
    highlightLength: number
}

interface State {
    blobLinesOrError?: string[] | ErrorLike
}

export class CodeExcerpt extends React.PureComponent<Props, State> {
    public state: State = {}
    private tableContainerElement: HTMLElement | null = null
    private propsChanges = new Subject<Props>()
    private visibilityChanges = new Subject<boolean>()
    private subscriptions = new Subscription()
    private visibilitySensorOffset = { bottom: -500 }

    public constructor(props: Props) {
        super(props)
        this.subscriptions.add(
            combineLatest(this.propsChanges, this.visibilityChanges)
                .pipe(
                    filter(([, isVisible]) => isVisible),
                    switchMap(([{ repoName, filePath, commitID, isLightTheme }]) =>
                        props.fetchHighlightedFileLines({
                            repoName,
                            commitID,
                            filePath,
                            isLightTheme,
                            disableTimeout: true,
                        })
                    ),
                    catchError(error => [asError(error)])
                )
                .subscribe(blobLinesOrError => {
                    this.setState({ blobLinesOrError })
                })
        )
    }

    public componentDidMount(): void {
        this.propsChanges.next(this.props)
    }

    public componentWillReceiveProps(nextProps: Props): void {
        this.propsChanges.next(nextProps)
    }

    public componentDidUpdate(prevProps: Props, prevState: State): void {
        if (this.tableContainerElement) {
            const visibleRows = this.tableContainerElement.querySelectorAll('table tr')
            for (const highlight of this.props.highlightRanges) {
                const code = visibleRows[highlight.line - this.getFirstLine()].lastChild as HTMLTableDataCellElement
                highlightNode(code, highlight.character, highlight.highlightLength)
            }
        }
    }

    public componentWillUnmount(): void {
        this.subscriptions.unsubscribe()
    }

    private getFirstLine(): number {
        return Math.max(0, Math.min(...this.props.highlightRanges.map(r => r.line)) - (this.props.context || 1))
    }

    private getLastLine(blobLines: string[] | undefined): number {
        const lastLine = Math.max(...this.props.highlightRanges.map(r => r.line)) + (this.props.context || 1)
        return blobLines ? Math.min(lastLine, blobLines.length) : lastLine
    }

    private onChangeVisibility = (isVisible: boolean): void => {
        this.visibilityChanges.next(isVisible)
    }

    public render(): JSX.Element | null {
        return (
            <VisibilitySensor
                onChange={this.onChangeVisibility}
                partialVisibility={true}
                offset={this.visibilitySensorOffset}
            >
                <code
                    className={`code-excerpt ${this.props.className || ''} ${
                        isErrorLike(this.state.blobLinesOrError) ? 'code-excerpt-error' : ''
                    }`}
                >
                    {this.state.blobLinesOrError && !isErrorLike(this.state.blobLinesOrError) && (
                        <div
                            ref={this.setTableContainerElement}
                            dangerouslySetInnerHTML={{ __html: this.makeTableHTML(this.state.blobLinesOrError) }}
                        />
                    )}
                    {this.state.blobLinesOrError && isErrorLike(this.state.blobLinesOrError) && (
                        <div
                            className="alert alert-danger "
                            dangerouslySetInnerHTML={{ __html: this.state.blobLinesOrError.message }}
                        />
                    )}
                    {!this.state.blobLinesOrError && (
                        <table>
                            <tbody>
                                {range(this.getFirstLine(), this.getLastLine(this.state.blobLinesOrError)).map(i => (
                                    <tr key={i}>
                                        <td className="line">{i + 1}</td>
                                        {/* create empty space to fill viewport (as if the blob content were already fetched, otherwise we'll overfetch) */}
                                        <td className="code"> </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    )}
                </code>
            </VisibilitySensor>
        )
    }

    private setTableContainerElement = (ref: HTMLElement | null) => {
        this.tableContainerElement = ref
    }

    private makeTableHTML(blobLines: string[]): string {
        return '<table>' + blobLines.slice(this.getFirstLine(), this.getLastLine(blobLines) + 1).join('') + '</table>'
    }
}
