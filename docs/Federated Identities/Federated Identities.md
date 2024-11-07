https://heimdalsecurity.com/blog/what-is-federated-identity/
https://www.okta.com/identity-101/what-is-federated-identity/
# [What is Federated Identity?](https://learning.oreilly.com/library/view/federated-identity-primer/9780124071896/)

## 2.2. Authentication and Authorization
>One of the key problems people have when studying authentication and authorization is understanding the difference between the two.

>However, when it comes to federated identity knowing the difference is very important. In fact, federated identity is based on the fact that the two concepts are not the same.
### 2.2.1 Authentication
>Authentication is broken down into two components: identification and verification. Identification occurs before verification. Identification is the process of stating who you are. This statement could be in the form of a username, an e-mail address, or some other method that identifies you.

>Verification is the process that a system goes through to check that you are indeed who you say you are. This is what most people think of when they think of authentication. They don’t realize that the first part of the process is establishing your identity. Verification can be performed in many ways. You supply a password, a Personal Identification Number (PIN), or use some type of biometric identifier.
#### 2.2.1.1.3 User Credential
>There is a certain type of digital certificate called a user certificate that is specifically designed for user authentication. After the certificate is created, the certificate is then mapped back to a user account. This user account is used to determine what access the user should have.

>When the user attempts to access a resource, a certificate will be requested.

>During the processing, the backend authentication system will look up the certificate and find the corresponding user account for that certificate.

>This information is then submitted to the resource. The resource will then make authorization decisions based on the user account.
#### 2.2.1.2.1. Mutual Authentication
>In mutual authentication, not only is the client authenticated, but the server is also authenticated. The server must do something to prove its identity. This could be in the form of a server certificate or some sort of private key. Once the server has been authenticated and the client trusts the server, then the client will send its credentials to the server. This provides for a more secure authentication process and a more secure environment overall.
### 2.2.2. Authorization
>Authorization is the process of specifying what a user is allowed to do.
## 2.4. Federated Service Model
>The federated service model can be broken up into two required components: the **IdP** and the **service provider**. Without both of these components, you do not have a federated identity solution. Each of these components has distinct characteristics and responsibilities.

>The IdP must trust the application or it will not send user information to the application. The application must trust the IdP or it will not trust the user identity information that comes from the IdP.
### 2.4.1. Identity Provider
