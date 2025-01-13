# First, we'll move my-chart the original article vendored in
# Git to the conventional location Holos uses to vendor charts.
mkdir -p components/my-chart/vendor/0.1.0
git mv my-chart components/my-chart/vendor/0.1.0/my-chart

# Helm value files move into the directory that will contain
# my-chart component definition.  components/my-chart is
# conventionally called the "my-chart component directory"
git mv my-values components/my-chart/my-values

# The config.json files are moved without changing their folder structure.
# We'll package the data up into an "environments config package" for reuse.
mkdir config
git mv appsets/4-final/env-config config/environments

# All of our components will reside in the components directory so
# the CUE files `holos init` produced may be moved to keep the
# repository root tidy.
mv *.cue components/

# The following files and directories from the original article and
# holos init are no longer relevant after the migration.
mkdir not-used
git mv appsets not-used/appsets
git mv example-apps not-used/example-apps
rm -f platform.metadata.json

# Make the commit
git add platform components cue.mod .gitignore
git commit -m 'reorganize to conventional holos layout'
