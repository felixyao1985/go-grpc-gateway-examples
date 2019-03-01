###JWT - GO

现成的库，不是官方，但是星高！有很多方法可以使用
 
    github.com/dgrijalva/jwt-go 
    
    
```
/*
测试Token
*/
Token,_ := jwt.GenerateToken("felix","123456")
fmt.Println(Token)
jwt.ParseToken(Token.Token)

```