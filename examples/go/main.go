package main

import (
	"github.com/pulumi/pulumi-statefulString/sdk/go/statefulString"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		myStatefulString, err := statefulString.NewStatefulString(ctx, "myStatefulString", &statefulString.StatefulStringArgs{
			String: pulumi.String("hello, world"),
			Triggers: pulumi.StringMap{
				"key": pulumi.String("value"),
			},
		})
		if err != nil {
			return err
		}
		ctx.Export("output", map[string]interface{}{
			"value": myStatefulString.String,
		})
		return nil
	})
}
