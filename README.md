DAWG
====

Fancy link expander workflow generator for [Alfred](https://www.alfredapp.com/).

## Usage

Put a json config with URL templates and substitutions to `~/.dawg.json`:

    {
      "www.datadoghq.com": {
        "keyword": "dog",
        "template": "https://app.datadoghq.com/screen/{id}",
        "substitutions": {
          "mydash":    { "id": "100" },
          "otherdash": { "id": "101" }
        }
      },
      "www.onelogin.com": {
        "keyword": "ol",
        "template": "https://app.onelogin.com/client/apps/select/{id}",
        "substitutions": {
          "datadog": { "id": "1000" },
          "google":  { "id": "1001" }
        }
      }
    }

Generate the workflow:

    dawg -generate

Install the workflow:

    open DAWG.alfredworkflow

## Screenshots!

![Alfred](https://github.com/v-yarotsky/dawg/blob/master/doc/screenshot.png?raw=true)

## License

[WTFPL](https://github.com/v-yarotsky/dawg/blob/master/LICENSE.txt?raw=true)

