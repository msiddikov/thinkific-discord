# thinkific-discord


Error when sending request to the server
panic: runtime error: invalid memory address or nil pointer dereference
[signal 0xc0000005 code=0x0 addr=0x40 pc=0x167a3e4]

goroutine 4236 [running]:
thinkific-discord/internal/discordBot.discordReq(0xc0004bae00)
        H:/ReactProjects/thinkific-discord/internal/discordBot/discordBot.go:52 +0x2a4
thinkific-discord/internal/discordBot.GetRoles()
        H:/ReactProjects/thinkific-discord/internal/discordBot/discordBot.go:85 +0x8a
thinkific-discord/internal/discordBot.UpdateRoles()
        H:/ReactProjects/thinkific-discord/internal/discordBot/discordBot.go:99 +0x19
reflect.Value.call({0x17ee240?, 0x193f870?, 0xc000099e28?}, {0x18e83b5, 0x4}, {0x1e96868, 0x0, 0xc0000fe198?})
        C:/Program Files/Go/src/reflect/value.go:556 +0x845
reflect.Value.Call({0x17ee240?, 0x193f870?, 0x1255d68?}, {0x1e96868, 0x0, 0x0})
        C:/Program Files/Go/src/reflect/value.go:339 +0xbf
github.com/go-co-op/gocron.callJobFuncWithParams({0x17ee240?, 0x193f870?}, {0x0, 0x0, 0xdf8475800?})
        C:/Users/NVidia/go/pkg/mod/github.com/go-co-op/gocron@v1.13.0/gocron.go:78 +0x135
github.com/go-co-op/gocron.(*executor).start.func1()
        C:/Users/NVidia/go/pkg/mod/github.com/go-co-op/gocron@v1.13.0/executor.go:81 +0x172
created by github.com/go-co-op/gocron.(*executor).start
        C:/Users/NVidia/go/pkg/mod/github.com/go-co-op/gocron@v1.13.0/executor.go:49 +0x71

C:\Users\Administrator\Desktop\thinkific-discord>pause