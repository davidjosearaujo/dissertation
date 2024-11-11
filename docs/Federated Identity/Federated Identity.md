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
>You can have an IdP set up internally, or you can use a service provider to provide identity services for you. Users, also called entities, will authenticate against the IdP’s credential store. The IdP will then allow access to user’s identity information. It’s important to note that an IdP does more than just authenticate a user. It also holds the user’s identity information. Upon authentication, this information can be sent to whichever trusted partner needs it.
#### 2.4.1.1. Credential Store
> The credential store, sometimes called the user store or the authentication store, is where the actual user credentials are stored.

>There are two main types of authentication stores being used with IdPs today: **databases and directory stores**.

>In general, with databases, credentials are stored in proprietary tables created by the user management application. One of the reasons **databases are often chosen as credential stores** is because a majority of developers have experience coding against a database, so it’s relatively easy for them to create code to authenticate users against one.
### 2.4.2. Service Provider
>Service providers are the entities that provide services to others. These services could be applications, infrastructure, or data services. As “the cloud” grows in popularity, many people have become aware of the three main cloud services models.
## 2.5 Federated Identity
>One of the biggest confusions that exist around federated identity is how it is related to federated authentication. 

>**Federated authentication** can be considered a subset of a federated identity solution. Your digital identity is basically who you are, what you do in the digital world, and other characteristics.

>Federated authentication is the actual login process that takes place. You log into one place and that login allows you access somewhere else. What happens after that login is where other components of federated identity may or may not kick in.

>Federated authentication does not necessarily require an IdP to be in place. There may be some other systems in place where information is passed from one system to another. As I mentioned before, an IdP is necessary for you to have a true federated identity solution. So, what does this mean? **This means that you can have federated authentication without federated identity.**

>Your information may only exist in one system. But, with federated identity, other systems can also have access this information.

>The key to federated identity is trust. The system that holds your information and the system that is requesting your information must trust each other.

>The system requesting the information has to trust the sender to ensure they are getting accurate and trustworthy information.

>The application does not perform any actions to verify the user’s identity itself. It just believes what the IdP says.

>Before an application will believe an IdP, a trust relationship must be established. The application must be configured with the address of the IdP that it will be trusting. The IdP must be configured with the address of the application. In most cases, some type of keys will be exchanged between the two entities to actually establish the relationship.
### 2.5.1 Authentication vs Authorization with Federated Identity
>**One of the key characteristics of federated identity is that authentication is abstracted from authorization.**

>With federated identity, the authentication request does not have to be performed by the application.

>You can use a third-party system, like an internal or external IdP to provide the authentication. The IdP sends the authentication information back to the application.

### 2.5.2 Federated Identity Advantages and Disadvantages
#### 2.5.2.1 Advantages
##### 2.5.2.1.7 Highly Extensible
>**Since the application itself does not know and generally does not care how the user was authenticated, you can use any method you choose for authentication.**