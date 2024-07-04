---
sidebar_position: 2
---

# Registration

Holos leverages a simple web app to collect and store platform attributes with a web form.  Register an account with the web app to create and retrieve the platform model.

```
holos register user
```

:::tip

Holos allows you to customize all of the sections and fields of your platform model.

:::


## Generate your Platform

Generate your platform configuration from the holos reference platform embedded in the `holos` executable.  Platform configuration is stored in a git repository.

```bash
mkdir holos-infra
cd holos-infra
holos generate platform holos
```

The generate command writes many files organized by platform component into the current directory

TODO: Put a table here describing key elements?

:::tip

Take a peek at `holos generate platform --help` to see other platforms embedded in the holos executable.

:::

## Push the Platform Form

```
holos push platform form .
```

## Fill in the form

TODO

## Pull the Platform Model

Once the platform model is saved, pull it into the holos-infra repository:

```
holos pull platform model .
```

## Render the Platform

With the platform model and the platform spec, you're ready to render the complete platform configuration:

```
holos render platform ./platform
```

## Summary
