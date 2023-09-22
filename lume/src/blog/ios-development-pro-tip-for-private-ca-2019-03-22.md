---
title: iOS Development Pro Tip for Private CA Usage
date: 2019-03-22
---

In iOS, in order to get HTTPS working with certs from a private CA; there's another step you need to do if your users are on iOS 10.3 or newer (statistically: yes this matters to you). In order to do this:

- Ensure they have installed the profile on their device
- Open Settings
- Select General
- Select Profiles
- Ensure your root CA name is visible in the profile list like this:

<style>
img {
  max-width: 400px;
}
</style>

![](/static/img/ios_profiles.png)

- Go up a level to General
- Select About
- Select Certificate Trust Settings
- Each root that has been installed via a profile will be listed below the heading Enable Full Trust For Root Certificates
- Users can toggle on/off trust for each root:

![](/static/img/ios_cert_trust.png)

Please understand that by doing this, users will potentially be vulnerable to a
[HTTPS man in the middle attack a-la Superfish](https://slate.com/technology/2015/02/lenovo-superfish-scandal-why-its-one-of-the-worst-consumer-computing-screw-ups-ever.html). Please ensure that you have appropriate measures in place to keep the signing key for the CA safe.

I hope this helps.
