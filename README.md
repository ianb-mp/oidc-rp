# OIDC test

A POC showing how to handle a Github Actions OIDC ID token as the RP (relying party)

As per [Github Docs](https://docs.github.com/en/actions/security-for-github-actions/security-hardening-your-deployments/configuring-openid-connect-in-cloud-providers#requesting-the-jwt-using-environment-variables) these steps are involved:

- Using workflow action runtime token to request a JWT from Github representing the OIDC ID token
- POST the JWT to our cloud provider's token endpoint
- receive an access token from the cloud provider

