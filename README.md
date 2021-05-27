# switcheroo
A mutating admission webhook that ensures all images are pulled from a private container registry.

It has been developed using [Kubebuilder](https://github.com/kubernetes-sigs/kubebuilder) which is a framework that reduces the amount of boilerplate code you have to write to implement a webhook.

The webhook works by modifying each image within a Pod specification. It targets containers and init containers
```spec.containers[*].image, spec.initContainers[*].image```

If you have a target registry host of ``foo.com`` then the replacement rules are as follows.
1. Images with no registry host get prefixed with the target registry host. ``hello-world:latest``  becomes ``foo.com/hello-world:latest``
2. Images with an existing registry host prefix get the host replaced with the target registry host. ``bar.com/hello-world:latest``  becomes ``foo.com/hello-world:latest``
2. Images already prefixed with the target registry host get left as is. ``foo.com/hello-world:latest``  becomes ``foo.com/hello-world:latest``