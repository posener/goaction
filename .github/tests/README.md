# Workflow communication test jobs

The workflow defined at /.github/workflows/testworkflow.yml runs the jobs defined in this directory.
For each one of them, it runs goaction to update to Action file in the project root directory:
`goaction -path ./.github/tests/<test name>`. This overrides the `Dockerfile` and
`action.yml` files which are then invoked by an action with `uses: ./`.

The idea is to test communication between two actions. Things that were set in one job or step are
correctly read in another job or step.