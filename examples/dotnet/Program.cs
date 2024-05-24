using System.Collections.Generic;
using System.Linq;
using Pulumi;
using statefulstring = Pulumi.statefulstring;

return await Deployment.RunAsync(() => 
{
    var myRandomResource = new statefulstring.Random("myRandomResource", new()
    {
        Length = 24,
    });

    return new Dictionary<string, object?>
    {
        ["output"] = 
        {
            { "value", myRandomResource.Result },
        },
    };
});

