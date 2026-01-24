# Our Security Policy

We take the security of Dragonfly very seriously. If you've discovered a security vulnerability, we appreciate your help
in disclosing it to us responsibly.

## Reporting a Vulnerability

To report a vulnerability in Dragonfly, please follow these steps:

1. Go to the **Security** tab in the repository on GitHub.
2. Click on the **Advisories** tab.
3. Click on **Report a vulnerability**.

You can expect a response from us within 24 hours to acknowledge that we've received your report. If you don't hear
back from us in that time, please reach out to any committer directly to make sure we got your message.

## Our Security Response Team

The Dragonfly maintainers act as our security response team. All security issues are diagnosed and discussed on a
private mailing list at [dragonfly-maintainers@googlegroups.com](mailto:dragonfly-maintainers@googlegroups.com) to
ensure confidentiality.

## Our Review Process

Once a committer confirms that the report is relevant, they will create a draft security advisory on GitHub. This is
where we'll discuss the issue with committers, the reporter(s), and Dragonfly's security advisors.

If you'd like to be part of this discussion, please provide your GitHub username so we can invite you. If you'd rather
not participate directly, you can ask to be kept updated via email.

If we accept the vulnerability, we'll work together to determine a timeline for developing a patch, publicly disclosing
the issue, and releasing the fix. If there's an embargo period before the public disclosure, we'll send an announcement
to our security mailing list, [dragonfly-developers@googlegroups.com](mailto:dragonfly-developers@googlegroups.com).
This announcement will outline the scope of the vulnerability, when the patch will be available, and the date of public
disclosure. We expect reporters to be part of the timeline discussion and to respect the agreed-upon dates for public
disclosure.

## Supported Versions

You can find information on which versions of Dragonfly are currently supported on our releases page. Security updates
may be provided for any Extended or Active release branch. If you discover a security issue in an older version, a
non-core package, or a dependency, please let our committers know using the same security mailing list you'd use to
report any other vulnerability.

## Joining the Security Announce Mailing List

If your organization or you as an individual use Dragonfly directly in a production or security-critical environment,
you are eligible to join our security announce mailing list. If you use Dragonfly through a vendor, we ask that you
have your vendor join the list instead.

To join, you'll need to be sponsored by a Dragonfly committer or security advisor and have a history of handling
non-public security information responsibly. Please don't request sponsorship through public channels, as the
membership of this list is not public.

### Mailing Lists

- **[dragonfly-developers@googlegroups.com](mailto:dragonfly-developers@googlegroups.com)**: Subscribe to receive
  security announcements.

### Confidentiality, Integrity, and Availability

We prioritize vulnerabilities that could compromise data confidentiality, allow for privilege elevation, or affect data
integrity. We also take availability issues, like Denial of Service (DoS) and resource exhaustion, very seriously.

We expect operators to configure their settings, role-based access control, and other features to create a hardened and
secure environment.
