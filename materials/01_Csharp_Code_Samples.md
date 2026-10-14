#### **RBAC**
Рассмотрим пример на C#, как может выглядеть доступ к методам на основе ролей.

```csharp
[Authorize(Roles="Cashier")]
public ActionResult SellItem(...)
{
    ...
}

[Authorize(Roles="Customer")]
public ActionResult AddItemToShoppingCart(...)
{
    ...
}
```

#### **CBAC**

Посмотрим пример на C#:

```csharp
[ClaimsAuthorize(Claims="Adult-older-than-18, Trust-worthy-client")]
public ActionResult RequestLoan(...)
{
    ...
}

[ClaimsAuthorize(Claims="Child")]
public ActionResult AddParent(...)
{
    ...
}
```

### Dependency Injection, IoC-контейнеры
Рассмотрим пример, как работает внедрение зависимостей на языке C#.

Для начала создаем интерфейс некоторого класса и его реализацию.

```csharp
public interface ISomeClass
{
    Task WriteMessage(string message);
}

public class SomeClass: ISomeClass
{
    private readonly ILogger<SomeClass> _logger;

    public SomeClass(ILogger<SomeClass> logger)
    {
        _logger = logger;
    }

    public Task WriteMessage(string message)
    {
        _logger.LogInformation("...");
        return Task.FromResult(0);
    }
}
```

Теперь зарегистрируем зависимость в файле Program.cs (.net6 +).

```csharp
...
builder.Services.AddTransient<ISomeClass, SomeClass>();
...
```

Наконец, можно прокинуть зависимость в желаемый класс.

```csharp
public class IndexModel : PageModel
{
    private readonly ISomeClass _someClass;

    public IndexModel(ISomeClass someClass)
    {
        _someClass = someClass;
    }

    public async Task OnGetAsync()
    {
        await _someClass.WriteMessage("Hooray!");
    }
}
```

Ура! Мы научились добавлять зависимости без явного вызова конструктора объектов.