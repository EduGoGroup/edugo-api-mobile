Este es un proyecto de golang como lo vistes en la documentacion, y claude.md
La teoria es que estoy usando viper para el tema de configuracion, tambien tengo carpetas de configuracion para editores de texto como .zed, .vscode .idea, etc etc, tambien tengo archivos make, y tambien tengo configurado algunas cosas en cicd, y docker file y docker compose

Como ves mi configuracion esta ampliable pero la siento desordenada y sobre todo configuraciones repetitiva o deprecada, o desactualizadas
En la carpeta @config estan yaml, donde en teoria gracias a viper dinamica segun el sufijo del ambiente, que se guarda en una variabke de memoria, aca creo que hay un poco de desorden, pero eso te dejo que lo analises tu, ya que se me aconsejo que no guardara claves en esos archivos sino en secretos, pero es alli donde empieza el problema

Si no puedo guardar datos como clave de postgre, o de rabbit, o de amazon, entonces empieza cristo a padecer, porque tengo que jugar con archivos env en local, porque una cosa es cuando estas desplegando en ambientes en servicios en la nube, como QA, produccion o demas, y otra es cuando estas trabajando de manera local, es alli donde empieza los problemas

Puedo tener de forma local tanto
* Desarrollo en un IDE, debo saber donde cual es el archivo de configuracon a buscar, podria decir que la memoria se guarde APP_ENV, y o asocio ejemplo en Jetbrain sea Intellig o Golang que tenga la variable con el valor local, la teoria diria que el codigo debe buscar mi archivo condig-local.yaml, pero como no puedo tener guardada las claves entonces debo tener un archivo .env, me parece engorroso, he visto sitios, que los secretos tienen una varibale y luego en el momento de correr se sustitutye, asi no tengo que andar creando cosas en memoria, y el programador solo debe indicar donde estan esos valores

* Desarrollo en Editor de texto como zed o fork de visual code, donde cada uno tiene su propia carpeta oculta, como .kiro, .zed, .vscode, etc, aca es muy parecido al punto anterior, deberia tener una forma de configurar que pueda correr mi api y pasarle cual es el archivo, mismo tema como donde guardo la memoria,

* Correr docker compose o dockerfile, para que la app pueda correr en contenedores, aca tuve que hacer una adaptacion para que el dockerfile pueda copiar una archivo de env_example mismo problema

* Ahhh y me faltaba Make, que es el mismo principio

Entonces necesito que hagas una revision de la metodologia implementada de como maneja y persiste las variables de configuracion, quiero ver si se aplica las mejores tecnicas, y simplificar el proceso para que , por un lado permita ser flexible, pero no complicado (como siento ahora que tiene mucha logica innecesaria)

En ese tema de la logica, en unos de los ultimos sprint se agrego una variable para guardar clave de amazon, esta bien, esa es la idea, pero no se aplicaron las mejores formas, y ahora resulta, que no se llena la variable por mas que la coloque en el env, ya que (y no entiendo porque deberia ser una variable que sustituya) no puedo colocarlo en config-xxx.yaml

Conclusion
* Analiza el sistema actual de manejo de variables de configuracion, como se extrae, como se inyecta para que lo usen en codigo, como es la jerarquia para trabajar multiple ambiente, como se comporta y debe ser para el uso local, tanto en IDE, editores de codigo, make y docker, como se puede hacer una implementacion que sea facil tanto en el consumo local, como en el momento que se vaya a desplegar en un futuro en una nube que usa secretos
* Crea un script aparte, que permita agregar nuevas variables en la configuracion, es decir si el dia de mañana necesito agregar una nueva variable como token de google, yo pueda decir clave, valor, y jerarquia (busca la manera mas facil, para indicar como puede ser la jerarqquia en yaml)
* El punto anterior, debe agregar dicha nueva variable, al estandar segun el archivo yaml, y segun el sistema que plantee para sustituir los secretos
Muestrame un plan de trabajo de como hacer el cambio, debes crear una rama de dev en el momento de la implementacion del proceso, y debes validar que pueda levantar el api, despues que se agregue una nueva variable de test

