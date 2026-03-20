# 功能说明
我们为Agent-Sandbox来生成一个后台管理UI，可以查看和管理当前的Sandbox实例，包括它们的状态、日志、资源使用情况等信息。这个UI将会提供一个直观的界面，帮助用户更好地理解和控制他们的Sandboxes。

## 技术要求：
1. UI使用英文；
2. 不需要默认进行Verification (End-to-End)操作；
3. 使用主流前端框架（如React、Vue）和daisyUI来构建UI；
4. 注意完全使用daisyUI的组件和样式，保持界面的一致性和美观性；
5. 左边有导航栏，右边是内容展示区，不必是传统的那种布局，可以有创新；

## 功能需求：
1. 结合../pkg/handler/handler.go实现Sandboxes管理功能，包括Sandboxes列表，创建Sandbox和删除Sandbox功能；
2. 结合../pkg/handler/handler.go实现Rest API for config管理功能，包括查看templates配置和修改templates配置功能；
3. 结合../pkg/handler/handler.go实现Rest API for pool management功能，包括Pool列表，删除Pool（相当于删除了此Pool下面全部的Sandboxes）和查看Pool详情，Pool详情包括属于Pool的Sandboxes列表，工具栏同样有删除此Pool的功能；
4. 实现日志查看功能，可以查看每个Sandbox的日志输出，pkg/handler/handler.go增加/logs API和handler实现，输出sandbox pod的日志内容，UI上可以通过点击某个Sandbox来查看它的日志输出，日志输出需要支持自动刷新；
5. 实现进入Sandbox Terminal功能，pkg/handler/handler.go增加/terminal API和handler实现，通过UI可以进入 Sandbox pod的terminal执行shell，整体设计参考日志(logs)；
6. 实现Sandbox Files管理功能，pkg/handler/handler.go增加/sandbox/files API和handler实现，通过UI进入Sandbox pod files管理，功能包括遍历文件、下载文件、上传文件等功能；