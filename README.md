# OIDC test

A POC showing how to handle a Github Actions OIDC ID token as the RP (relying party)

## GH Action

As per [Github Docs](https://docs.github.com/en/actions/security-for-github-actions/security-hardening-your-deployments/configuring-openid-connect-in-cloud-providers#requesting-the-jwt-using-environment-variables) these steps are involved:

- Using workflow action runtime token to request a JWT from Github representing the OIDC ID token
- POST the JWT to our cloud provider's token endpoint
- receive an access token from the cloud provider

## Cloud provider

The `oidc-rp` Go binary runs on the cloud provider side and receives the Github ID token & verifies it. It is upto the cloud provider to provide a way for the user to define a trust relationship with Github i.e. this is outside the scope of OIDC. Conceptually, this could work like this:

- user logs into customer portal
- in account settings, under 'SSO' there is an option to add an OIDC identiy provider. At a minimum, these values are required:
  - provider url e.g. `https://token.actions.githubusercontent.com`
  - client ID e.g. `https://github.com/ianb-mp`
- the endpoint that handles the Github ID token can lookup these values to perform the verification against
- if verification succeeds, then temporary API credentials are generated and returned to the client

