# Holos caches charts here
mkdir -p components/my-chart/vendor/0.1.0
git mv my-chart components/my-chart/vendor/0.1.0/my-chart
# Values move to the helm chart component directory
git mv my-values components/my-chart/my-values
# config.json files with environment attributes
mkdir config
git mv appsets/4-final/env-config config/environments

# Holos CUE definitions
mv *.cue components/

# Not part of the migration
mkdir not-used
git mv appsets not-used/appsets
git mv example-apps not-used/example-apps
rm -f platform.metadata.json

# Make the commit
git add platform components cue.mod .gitignore
git commit -m 'reorganize to conventional holos layout'
