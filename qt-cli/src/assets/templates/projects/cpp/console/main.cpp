#include <QCoreApplication>
{{- if .useTranslation }}
#include <QLocale>
#include <QTranslator>
{{- end }}

int main(int argc, char *argv[])
{
    QCoreApplication app(argc, argv);
{{ if .useTranslation }}
    QTranslator translator;
    const QStringList uiLanguages = QLocale::system().uiLanguages();
    for (const QString &locale : uiLanguages) {
        const QString baseName = "{{ .name }}_" + QLocale(locale).name();
        if (translator.load(":/i18n/" + baseName)) {
            app.installTranslator(&translator);
            break;
        }
    }
{{- end }}
    // Set up code that uses the Qt event loop here.
    // Call app.quit() or app.exit() to quit the application.
    // A not very useful example would be including
    // #include <QTimer>
    // near the top of the file and calling
    // QTimer::singleShot(5000, &a, &QCoreApplication::quit);
    // which quits the application after 5 seconds.

    // If you do not need a running Qt event loop, remove the call
    // to app.exec() or use the Non-Qt Plain C++ Application template.

    return app.exec();
}
