Holos Platform

| Folder     | Description                                                                                             |
| -          | -                                                                                                       |
| forms      | Contains Platform and Project form and model definitions                                                |
| platform   | Contains the Platform resource that defines how to render the configuration for all Platform Components |
| components | Contains BuildPlan resources which define how to render individual Platform Components                  |

## Forms

Populate the platform web form from the cue configuration.

```bash
holos push platform form .
```

Fill out the form to build the Platform Model.

## Platform Model

Pull the most recent Platform Model each time the platform or components are
rendered.

```
holos pull platform config .
```
