import * as pulumi from "@pulumi/pulumi";
import * as statefulstring from "@pulumi/statefulstring";
// import * as statefulstring from "../nodejs/bin/StatefulString";

const mystring = new statefulstring.StatefulString("mystatefulstring", {
    string: "Hello, World!333", triggers: {
        "foo": "bar2",
    }
});

export const output = {
    value: mystring.string,
};
