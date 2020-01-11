Для запуска сервера перейдите в директорию из поисковой системой и выполните команду ./server

<div class="container" style="margin-bottom: 100px; display: block;">
    <ul>
      <li>
        <h3>Добавление и обновление индекса</h3>
        <pre>data = {
    "id": { ID вашей записи },
    "text": json.dumps({"field__indexing": "{ Пример текста }"}),
    "table-name": "{ Имя таблицы }",
}
response = requests.post("http://127.0.0.1:8000/add-index/", data=data)
        </pre>
        <h5>Для обновление записи в базе данных отправте запрос с именем таблицы и апи ключом там где вы хотите обновить запис</h5>
        <p>В качестве { Имени таблицы } используйте только латиницу и по желанию можно использовать только знак "_"</p>
        <p>Для обработки поля, после его название используйте __indexing</p>
        <p>Все поля без обозначение __indexing будут сохранены без обработки</p>
        <p>Крайне не рекомендуется индексировать разные поля в одной таблице</p>
      </li>
      <li>
        <h3>Удаление индекса</h3>
        <pre>data = {
    "id": { Id записи которую хотите удалить },
    "table-name": "{ Имя таблицы }"
}
requests.post("http://127.0.0.1:8000/remove-index/", data=data)
        </pre>
      </li>
      <li>
        <h3>Поиск</h3>
        <pre>{
    "text": json.dumps({"field__search": "{ Поисковий запрос }"}),
    "table-name": "{ Имя таблицы }",
}
requests.get("http://127.0.0.1:8000/search/", params=data)
        </pre>
        <p>Чтобы определить поле по которому нужно искать используйте постфикс __search</p>
        <p>Поиск будет работать только в зарание проиндексированых полях, там где вы использовали __indexing при добавление или обновление индекса</p>
        <p>В качества ответа сервер возвращает масив словарей из ID записи и Score ошыбкой</p>
      </li>
    </ul>
  </div>
