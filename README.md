DAWG
====

Fancy link expander for use with [Alfred](https://www.alfredapp.com/).

## Usage

Put a json config with URL templates and substitutions:

    {
      "datadog": {
        "template": "https://app.datadoghq.com/screen/{id}",
        "substitutions": {
          "mydash":    { "id": "100" },
          "otherdash": { "id": "101" }
        }
      },
      "onelogin": {
        "template": "https://app.onelogin.com/client/apps/select/{id}",
        "substitutions": {
          "datadog": { "id": "1000" },
          "google":  { "id": "1001" }
        }
      }
    }

When you run this:

    dawg -s datadog

You will get this:

    <?xml version="1.0" encoding="UTF-8"?>
    <items>
      <item uid="dawg:mydash" autocomplete="mydash">
        <title>mydash</title>
        <arg>https://app.datadoghq.com/screen/100</arg>
      </item>
      <item uid="dawg:otherdash" autocomplete="otherdash">
        <title>otherdash</title>
        <arg>https://app.datadoghq.com/screen/101</arg>
      </item>
    </items>

And when you run this:

    dawg -s datadog my

You will get this:

    <?xml version="1.0" encoding="UTF-8"?>
    <items>
      <item uid="dawg:mydash" autocomplete="mydash">
        <title>mydash</title>
        <arg>https://app.datadoghq.com/screen/100</arg>
      </item>
    </items>

## Screenshots!

Here's how a simple "script filter + open URL" Alfred workflow looks like:

![Alfred](https://github.com/v-yarotsky/dawg/blob/master/doc/screenshot.png?raw=true)
![Alfred Workflow Setup](https://github.com/v-yarotsky/dawg/blob/master/doc/workflow.png?raw=true)

## License

[WTFPL](https://github.com/v-yarotsky/dawg/blob/master/LICENSE.txt?raw=true)

