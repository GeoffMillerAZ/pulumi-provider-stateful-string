import pulumi
import pulumi_statefulstring as statefulstring

my_random_resource = statefulstring.Random("myRandomResource", length=24)
pulumi.export("output", {
    "value": my_random_resource.result,
})
