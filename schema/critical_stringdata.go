// Code generated by stringdata. DO NOT EDIT.

package schema

// CriticalSchemaJSON is the content of the file "critical.schema.json".
const CriticalSchemaJSON = `{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "$id": "critical.schema.json#",
  "title": "Critical configuration",
  "description": "Critical configuration for a Sourcegraph site.",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "auth.userOrgMap": {
      "description": "Ensure that matching users are members of the specified orgs (auto-joining users to the orgs if they are not already a member). Provide a JSON object of the form ` + "`" + `{\"*\": [\"org1\", \"org2\"]}` + "`" + `, where org1 and org2 are orgs that all users are automatically joined to. Currently the only supported key is ` + "`" + `\"*\"` + "`" + `.",
      "type": "object",
      "additionalProperties": {
        "type": "array",
        "items": {
          "type": "string"
        }
      },
      "examples": [{ "*": ["myorg1"] }],
      "hide": true
    },
    "log": {
      "description": "Configuration for logging and alerting, including to external services.",
      "type": "object",
      "additionalProperties": false,
      "properties": {
        "sentry": {
          "description": "Configuration for Sentry",
          "type": "object",
          "additionalProperties": false,
          "properties": {
            "dsn": {
              "description": "Sentry Data Source Name (DSN). Per the Sentry docs (https://docs.sentry.io/quickstart/#about-the-dsn), it should match the following pattern: '{PROTOCOL}://{PUBLIC_KEY}@{HOST}/{PATH}{PROJECT_ID}'.",
              "type": "string",
              "pattern": "^https?://"
            }
          }
        }
      },
      "examples": [{ "sentry": { "dsn": "https://mykey@sentry.io/myproject" } }],
      "group": "Misc."
    },
    "externalURL": {
      "description": "The externally accessible URL for Sourcegraph (i.e., what you type into your browser). Previously called ` + "`" + `appURL` + "`" + `.",
      "type": "string",
      "examples": ["https://sourcegraph.example.com"]
    },
    "lightstepAccessToken": {
      "description": "Access token for sending traces to LightStep.",
      "type": "string",
      "group": "Misc."
    },
    "lightstepProject": {
      "description": "The project ID on LightStep that corresponds to the ` + "`" + `lightstepAccessToken` + "`" + `, only for generating links to traces. For example, if ` + "`" + `lightstepProject` + "`" + ` is ` + "`" + `mycompany-prod` + "`" + `, all HTTP responses from Sourcegraph will include an X-Trace header with the URL to the trace on LightStep, of the form ` + "`" + `https://app.lightstep.com/mycompany-prod/trace?span_guid=...&at_micros=...` + "`" + `.",
      "type": "string",
      "examples": ["myproject"],
      "group": "Misc."
    },
    "useJaeger": {
      "description": "Use local Jaeger instance for tracing. Kubernetes cluster deployments only.\n\nAfter enabling Jaeger and updating your Kubernetes cluster, ` + "`" + `kubectl get pods` + "`" + `\nshould display pods prefixed with ` + "`" + `jaeger-cassandra` + "`" + `,\n` + "`" + `jaeger-collector` + "`" + `, and ` + "`" + `jaeger-query` + "`" + `. ` + "`" + `jaeger-collector` + "`" + ` will start\ncrashing until you initialize the Cassandra DB. To do so, do the\nfollowing:\n\n1. Install [` + "`" + `cqlsh` + "`" + `](https://pypi.python.org/pypi/cqlsh).\n1. ` + "`" + `kubectl port-forward $(kubectl get pods | grep jaeger-cassandra | awk '{ print $1 }') 9042` + "`" + `\n1. ` + "`" + `git clone https://github.com/uber/jaeger && cd jaeger && MODE=test ./plugin/storage/cassandra/schema/create.sh | cqlsh` + "`" + `\n1. ` + "`" + `kubectl port-forward $(kubectl get pods | grep jaeger-query | awk '{ print $1 }') 16686` + "`" + `\n1. Go to http://localhost:16686 to view the Jaeger dashboard.",
      "type": "boolean",
      "group": "Misc."
    },
    "htmlHeadTop": {
      "description": "HTML to inject at the top of the ` + "`" + `<head>` + "`" + ` element on each page, for analytics scripts",
      "type": "string",
      "group": "Misc."
    },
    "htmlHeadBottom": {
      "description": "HTML to inject at the bottom of the ` + "`" + `<head>` + "`" + ` element on each page, for analytics scripts",
      "type": "string",
      "group": "Misc."
    },
    "htmlBodyTop": {
      "description": "HTML to inject at the top of the ` + "`" + `<body>` + "`" + ` element on each page, for analytics scripts",
      "type": "string",
      "group": "Misc."
    },
    "htmlBodyBottom": {
      "description": "HTML to inject at the bottom of the ` + "`" + `<body>` + "`" + ` element on each page, for analytics scripts",
      "type": "string",
      "group": "Misc."
    },
    "licenseKey": {
      "description": "The license key associated with a Sourcegraph product subscription, which is necessary to activate Sourcegraph Enterprise functionality. To obtain this value, contact Sourcegraph to purchase a subscription.",
      "type": "string",
      "group": "Sourcegraph Enterprise license"
    },
    "auth.providers": {
      "description": "The authentication providers to use for identifying and signing in users. See instructions below for configuring SAML, OpenID Connect (including G Suite), and HTTP authentication proxies. Multiple authentication providers are supported (by specifying multiple elements in this array).",
      "type": "array",
      "items": {
        "required": ["type"],
        "properties": {
          "type": {
            "type": "string",
            "enum": ["builtin", "saml", "openidconnect", "http-header", "github", "gitlab"]
          }
        },
        "oneOf": [
          { "$ref": "#/definitions/BuiltinAuthProvider" },
          { "$ref": "#/definitions/SAMLAuthProvider" },
          { "$ref": "#/definitions/OpenIDConnectAuthProvider" },
          { "$ref": "#/definitions/HTTPHeaderAuthProvider" },
          { "$ref": "#/definitions/GitHubAuthProvider" },
          { "$ref": "#/definitions/GitLabAuthProvider" }
        ],
        "!go": {
          "taggedUnionType": true
        }
      },
      "group": "Authentication",
      "default": [{ "type": "builtin", "allowSignup": true }]
    },
    "auth.public": {
      "description": "Allows anonymous visitors full read access to repositories, code files, search, and other data (except site configuration).\n\nSECURITY WARNING: If you enable this, you must ensure that only authorized users can access the server (using firewall rules or an external proxy, for example).\n\nRequires usage of the builtin authentication provider.",
      "type": "boolean",
      "default": false,
      "group": "Authentication"
    },
    "auth.sessionExpiry": {
      "type": "string",
      "description": "The duration of a user session, after which it expires and the user is required to re-authenticate. The default is 90 days. There is typically no need to set this, but some users may have specific internal security requirements.\n\nThe string format is that of the Duration type in the Go time package (https://golang.org/pkg/time/#ParseDuration). E.g., \"720h\", \"43200m\", \"2592000s\" all indicate a timespan of 30 days.\n\nNote: changing this field does not affect the expiration of existing sessions. If you would like to enforce this limit for existing sessions, you must log out currently signed-in users. You can force this by removing all keys beginning with \"session_\" from the Redis store:\n\n* For deployments using ` + "`" + `sourcegraph/server` + "`" + `: ` + "`" + `docker exec $CONTAINER_ID redis-cli --raw keys 'session_*' | xargs docker exec $CONTAINER_ID redis-cli del` + "`" + `\n* For cluster deployments: \n  ` + "`" + `` + "`" + `` + "`" + `\n  REDIS_POD=\"$(kubectl get pods -l app=redis-store -o jsonpath={.items[0].metadata.name})\";\n  kubectl exec \"$REDIS_POD\" -- redis-cli --raw keys 'session_*' | xargs kubectl exec \"$REDIS_POD\" -- redis-cli --raw del;\n  ` + "`" + `` + "`" + `` + "`" + `\n",
      "default": "2160h",
      "examples": ["168h"],
      "group": "Authentication"
    },
    "auth.disableUsernameChanges": {
      "description": "Prevent users from changing their username after account creation.",
      "type": "boolean",
      "default": false
    },
    "update.channel": {
      "description": "The channel on which to automatically check for Sourcegraph updates.",
      "type": ["string"],
      "enum": ["release", "none"],
      "default": "release",
      "examples": ["none"],
      "group": "Misc."
    }
  },
  "definitions": {
    "BuiltinAuthProvider": {
      "description": "Configures the builtin username-password authentication provider.",
      "type": "object",
      "additionalProperties": false,
      "required": ["type"],
      "properties": {
        "type": {
          "type": "string",
          "const": "builtin"
        },
        "allowSignup": {
          "description": "Allows new visitors to sign up for accounts. The sign-up page will be enabled and accessible to all visitors.\n\nSECURITY: If the site has no users (i.e., during initial setup), it will always allow the first user to sign up and become site admin **without any approval** (first user to sign up becomes the admin).",
          "type": "boolean",
          "default": false
        }
      }
    },
    "OpenIDConnectAuthProvider": {
      "description": "Configures the OpenID Connect authentication provider for SSO.",
      "type": "object",
      "additionalProperties": false,
      "required": ["type", "issuer", "clientID", "clientSecret"],
      "properties": {
        "type": {
          "type": "string",
          "const": "openidconnect"
        },
        "displayName": { "$ref": "#/definitions/AuthProviderCommon/properties/displayName" },
        "configID": {
          "description": "An identifier that can be used to reference this authentication provider in other parts of the config. For example, in configuration for a code host, you may want to designate this authentication provider as the identity provider for the code host.",
          "type": "string"
        },
        "issuer": {
          "description": "The URL of the OpenID Connect issuer.\n\nFor Google Apps: https://accounts.google.com",
          "type": "string",
          "format": "uri",
          "pattern": "^https?://"
        },
        "clientID": {
          "description": "The client ID for the OpenID Connect client for this site.\n\nFor Google Apps: obtain this value from the API console (https://console.developers.google.com), as described at https://developers.google.com/identity/protocols/OpenIDConnect#getcredentials",
          "type": "string",
          "pattern": "^[^<]"
        },
        "clientSecret": {
          "description": "The client secret for the OpenID Connect client for this site.\n\nFor Google Apps: obtain this value from the API console (https://console.developers.google.com), as described at https://developers.google.com/identity/protocols/OpenIDConnect#getcredentials",
          "type": "string",
          "pattern": "^[^<]"
        },
        "requireEmailDomain": {
          "description": "Only allow users to authenticate if their email domain is equal to this value (example: mycompany.com). Do not include a leading \"@\". If not set, all users on this OpenID Connect provider can authenticate to Sourcegraph.",
          "type": "string",
          "pattern": "^[^<@]"
        }
      }
    },
    "SAMLAuthProvider": {
      "description": "Configures the SAML authentication provider for SSO.\n\nNote: if you are using IdP-initiated login, you must have *at most one* SAMLAuthProvider in the ` + "`" + `auth.providers` + "`" + ` array.",
      "type": "object",
      "additionalProperties": false,
      "required": ["type"],
      "dependencies": {
        "serviceProviderCertificate": ["serviceProviderPrivateKey"],
        "serviceProviderPrivateKey": ["serviceProviderCertificate"],
        "signRequests": ["serviceProviderCertificate", "serviceProviderPrivateKey"]
      },
      "properties": {
        "type": {
          "type": "string",
          "const": "saml"
        },
        "configID": {
          "description": "An identifier that can be used to reference this authentication provider in other parts of the config. For example, in configuration for a code host, you may want to designate this authentication provider as the identity provider for the code host.",
          "type": "string"
        },
        "displayName": { "$ref": "#/definitions/AuthProviderCommon/properties/displayName" },
        "serviceProviderIssuer": {
          "description": "The name of this SAML Service Provider, which is used by the Identity Provider to identify this Service Provider. It defaults to https://sourcegraph.example.com/.auth/saml/metadata (where https://sourcegraph.example.com is replaced with this Sourcegraph instance's \"externalURL\"). It is only necessary to explicitly set the issuer if you are using multiple SAML authentication providers.",
          "type": "string"
        },
        "identityProviderMetadataURL": {
          "description": "The SAML Identity Provider metadata URL (for dynamic configuration of the SAML Service Provider).",
          "type": "string",
          "format": "uri",
          "pattern": "^https?://"
        },
        "identityProviderMetadata": {
          "description": "The SAML Identity Provider metadata XML contents (for static configuration of the SAML Service Provider). The value of this field should be an XML document whose root element is ` + "`" + `<EntityDescriptor>` + "`" + ` or ` + "`" + `<EntityDescriptors>` + "`" + `.",
          "type": "string"
        },
        "serviceProviderCertificate": {
          "description": "The SAML Service Provider certificate in X.509 encoding (begins with \"-----BEGIN CERTIFICATE-----\"). This certificate is used by the Identity Provider to validate the Service Provider's AuthnRequests and LogoutRequests. It corresponds to the Service Provider's private key (` + "`" + `serviceProviderPrivateKey` + "`" + `).",
          "type": "string",
          "$comment": "The pattern matches either X.509 encoding or an env var.",
          "pattern": "^(-----BEGIN CERTIFICATE-----\n|\\$)",
          "minLength": 1
        },
        "serviceProviderPrivateKey": {
          "description": "The SAML Service Provider private key in PKCS#8 encoding (begins with \"-----BEGIN PRIVATE KEY-----\"). This private key is used to sign AuthnRequests and LogoutRequests. It corresponds to the Service Provider's certificate (` + "`" + `serviceProviderCertificate` + "`" + `).",
          "type": "string",
          "$comment": "The pattern matches either PKCS#8 encoding or an env var.",
          "pattern": "^(-----BEGIN PRIVATE KEY-----\n|\\$)",
          "minLength": 1
        },
        "nameIDFormat": {
          "description": "The SAML NameID format to use when performing user authentication.",
          "type": "string",
          "pattern": "^urn:",
          "default": "urn:oasis:names:tc:SAML:2.0:nameid-format:persistent",
          "examples": [
            "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress",
            "urn:oasis:names:tc:SAML:1.1:nameid-format:persistent",
            "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified",
            "urn:oasis:names:tc:SAML:2.0:nameid-format:emailAddress",
            "urn:oasis:names:tc:SAML:2.0:nameid-format:persistent",
            "urn:oasis:names:tc:SAML:2.0:nameid-format:transient",
            "urn:oasis:names:tc:SAML:2.0:nameid-format:unspecified"
          ]
        },
        "signRequests": {
          "description": "Sign AuthnRequests and LogoutRequests sent to the Identity Provider using the Service Provider's private key (` + "`" + `serviceProviderPrivateKey` + "`" + `). It defaults to true if the ` + "`" + `serviceProviderPrivateKey` + "`" + ` and ` + "`" + `serviceProviderCertificate` + "`" + ` are set, and false otherwise.",
          "type": "boolean",
          "!go": { "pointer": true }
        },
        "insecureSkipAssertionSignatureValidation": {
          "description": "Whether the Service Provider should (insecurely) accept assertions from the Identity Provider without a valid signature.",
          "type": "boolean",
          "default": false
        }
      }
    },
    "HTTPHeaderAuthProvider": {
      "description": "Configures the HTTP header authentication provider (which authenticates users by consulting an HTTP request header set by an authentication proxy such as https://github.com/bitly/oauth2_proxy).",
      "type": "object",
      "additionalProperties": false,
      "required": ["type", "usernameHeader"],
      "properties": {
        "type": {
          "type": "string",
          "const": "http-header"
        },
        "usernameHeader": {
          "description": "The name (case-insensitive) of an HTTP header whose value is taken to be the username of the client requesting the page. Set this value when using an HTTP proxy that authenticates requests, and you don't want the extra configurability of the other authentication methods.",
          "type": "string",
          "examples": ["X-Forwarded-User"]
        },
        "stripUsernameHeaderPrefix": {
          "description": "The prefix that precedes the username portion of the HTTP header specified in ` + "`" + `usernameHeader` + "`" + `. If specified, the prefix will be stripped from the header value and the remainder will be used as the username. For example, if using Google Identity-Aware Proxy (IAP) with Google Sign-In, set this value to ` + "`" + `accounts.google.com:` + "`" + `.",
          "type": "string",
          "examples": ["accounts.google.com:"]
        }
      }
    },
    "GitHubAuthProvider": {
      "description": "Configures the GitHub (or GitHub Enterprise) OAuth authentication provider for SSO. In addition to specifying this configuration object, you must also create a OAuth App on your GitHub instance: https://developer.github.com/apps/building-oauth-apps/creating-an-oauth-app/. When a user signs into Sourcegraph or links their GitHub account to their existing Sourcegraph account, GitHub will prompt the user for the repo scope.",
      "type": "object",
      "additionalProperties": false,
      "required": ["type", "clientID", "clientSecret"],
      "properties": {
        "type": {
          "type": "string",
          "const": "github"
        },
        "url": {
          "type": "string",
          "description": "URL of the GitHub instance, such as https://github.com or https://github-enterprise.example.com.",
          "default": "https://github.com/"
        },
        "clientID": {
          "type": "string",
          "description": "The Client ID of the GitHub OAuth app, accessible from https://github.com/settings/developers (or the same path on GitHub Enterprise)."
        },
        "clientSecret": {
          "type": "string",
          "description": "The Client Secret of the GitHub OAuth app, accessible from https://github.com/settings/developers (or the same path on GitHub Enterprise)."
        },
        "displayName": { "$ref": "#/definitions/AuthProviderCommon/properties/displayName" },
        "allowSignup": {
          "description": "Allows new visitors to sign up for accounts via GitHub authentication. If false, users signing in via GitHub must have an existing Sourcegraph account, which will be linked to their GitHub identity after sign-in.",
          "default": false,
          "type": "boolean"
        }
      }
    },
    "GitLabAuthProvider": {
      "description": "Configures the GitLab OAuth authentication provider for SSO. In addition to specifying this configuration object, you must also create a OAuth App on your GitLab instance: https://docs.gitlab.com/ee/integration/oauth_provider.html. The application should have ` + "`" + `api` + "`" + ` and ` + "`" + `read_user` + "`" + ` scopes and the callback URL set to the concatenation of your Sourcegraph instance URL and \"/.auth/gitlab/callback\".",
      "type": "object",
      "additionalProperties": false,
      "required": ["type", "clientID", "clientSecret"],
      "properties": {
        "type": {
          "type": "string",
          "const": "gitlab"
        },
        "url": {
          "type": "string",
          "description": "URL of the GitLab instance, such as https://gitlab.com or https://gitlab.example.com.",
          "default": "https://gitlab.com/"
        },
        "clientID": {
          "type": "string",
          "description": "The Client ID of the GitLab OAuth app, accessible from https://gitlab.com/oauth/applications (or the same path on your private GitLab instance)."
        },
        "clientSecret": {
          "type": "string",
          "description": "The Client Secret of the GitLab OAuth app, accessible from https://gitlab.com/oauth/applications (or the same path on your private GitLab instance)."
        },
        "displayName": { "$ref": "#/definitions/AuthProviderCommon/properties/displayName" }
      }
    },
    "AuthProviderCommon": {
      "$comment": "This schema is not used directly. The *AuthProvider schemas refer to its properties directly.",
      "description": "Common properties for authentication providers.",
      "type": "object",
      "properties": {
        "displayName": {
          "description": "The name to use when displaying this authentication provider in the UI. Defaults to an auto-generated name with the type of authentication provider and other relevant identifiers (such as a hostname).",
          "type": "string"
        }
      }
    }
  }
}
`
