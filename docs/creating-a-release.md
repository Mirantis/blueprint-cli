# Creating a release

The release process is automated using github actions that trigger when a release is created on the github page.

1. Open releases on the github page
2. Create a pre-release which includes
  a. A tag for the latest commit on main. Use semantic versioning: `X.Y.Z`
  b. The auto generated changelog
  c. Check the pre-release box
  d. Publish the release
3. CI will trigger and begin the release process
  a. Run through all tests (lint, unit, integration)
  b. Build the release binaries
  c. Publish the binaries to the following repo's release pages
    i. https://github.com/mirantiscontainers/blueprint/releases
    i. https://github.com/mirantiscontainers/blueprint-cli/releases
4. Once CI finished, take a look at the binaries and make sure they look good
5. Change the release from pre-release to latest on the github page
