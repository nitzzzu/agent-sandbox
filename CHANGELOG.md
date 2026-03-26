# CHANGELOG


V0.1.0 - 2026-01-06
- Initial release of the project.
-------------------------- 
V0.1.1 - 2026-01-07
- Add sandbox template startupProbe config, fix get instance ip failed when it not ready in sometimes.
-------------------------- 
V0.2.0 - 2026-01-27
- Support E2B Protocol with SDK compatibility.
- Add E2B Code-interpreter support.
- Add E2B Desktop support with VNC and GUI applications.
- Support scale-down by timeout mechanism.
-------------------------- 
V0.3.0 - 2026-03-03
- Add Sandbox template Pool feature, which can pre-create sandbox instances for faster allocation.
- Add dynamic Sandbox template, which can create sandbox instances with template by pattern.
- Add OpenClaw template.
- Fix get default port bug.
- Fix httpServer WriteTimeout config bug.
-------------------------- 
V0.3.1 - 2026-03-12
- Add dynamic templates config load from configmap.
- Template pool support warmup feature, which can pre-run some commands or scripts to keep the sandbox instance warmup and low  resource consumption.
- Add skills for agent use.
-------------------------- 
V0.4.0 - 2026-03-19
- Add UI for sandbox management, which can view sandbox instances, templates, and logs.
-------------------------- 
V0.4.1 - 2026-03-26
- Add sandboxes events and metrics, which can monitor the status and performance of sandbox instances.
- Add template resources config for sandbox instance.
- Add template metadata config, which can customize sandbox-template specific config for different use cases with go-template format, e.g. #5.
- Add UI for sandbox-template config, which can view and edit the sandbox-template(ReplicasSet) config in UI.
-------------------------- 
