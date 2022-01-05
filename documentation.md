# Abstract
Voici la documentation du projet, ici nous considérons que chaque partie est indépendante des autres. C'est à vous de remettre le projet et les containers dans l'état initial afin que les commandes pour chaque partie fonctionnent.

# partie 1 serveur statique avec apache httpd
Ici, nous devons simplement faire un serveur web qui sert du contenu statique.
Déplacons nous dans le dossier step1, construisons l'image et lançons le container.

```bash
cd step1
sudo docker build -t static .
sudo docker run -d -p 8080:80 static
```
Pour voir le contenu static, vous pouvez aller avec un navigateur web sur localhost:8080

Pour vérifier que ce sont bien les fichiers, vous pouvez comparer entre le dossier landing-page et le dossier à l'intérieur du container de la manière suivante:
```bash
sudo docker exec -it <name/ID> /bin/sh
ls /usr/local/apache2/htdocs/
```
verifier que les fichiers (js,html,css) sont bien présents et identiques que ceux du dossier landing-page.
## configuration dans step1/:
- Dockerfile qui copie les fichiers statiques du dossier landing-page dans le dossier /usr/local/apache2/htdocs/
- Dossier `landing-page` qui contient les fichiers html, css, js a afficher par le serveur web.


# Partie 2 serveur dynamique
Dans le dossier step2, nous avons un dockerfile qui copie les fichiers go.* et main.go, compile le code et lance le serveur.

Le code génère une liste de personne (nom, prénom) ayant des hobbys, des interets et des capacités générées aléatoirement.
Le code est écrit en [Go](https://go.dev/) et utilise le framework [gin]("github.com/gin-gonic/gin") pour gérer le routage ainsi que la génération des réponses.

Pour tester: 
SUPPRIMER le container de la partie 1 (afin de pas avoir de conflit de port)
puis Builder l'image à partir du dockerfile
```bash
cd step2
sudo docker build -t dynamic .
```
puis lancer le container
```bash
sudo docker run -d -p 8080:80 dynamic
sudo docker ps
````
Lancez un navigateur et aller sur localhost:8080, vérifiez que des données JSON vous sont bien retournées. Les données sont aléatoires il n'y a pas de contrôle sur le contenu si ce n'est la constante `MAX_PERSON` dans `main.go` qui gère le nombre max de personne dans la réponse. Les différentes valeurs inclues dans les données sont définis dans le fichier `main.go` lignes 20-24 si vous voulez en ajouter/changer.

# Partie 3, reverse proxy
Comme proposé en cours, nous avons choisis d'utilisé Nginx pour gérer le reverse proxy.
Les fichiers sont dans le dossier step3. Afin de simplifier le nombre de commande, nous avons utilisé docker-compose qui lance les étapes 1,2 et 3.  
Les étapes une et deux ne sont pas modifiées, l'étape 3 consistes juste en l'ajout du reverse proxy. La configuration de ce dernier se fait dans le fichier nginx.conf dans le dossier step3.
les lignes modifiées par rapport à la configuration de nginx.conf par défaut sont :
````
#include /etc/nginx/conf.d/*.conf;
	server{
		listen 80;
		location /dyn/ {
			proxy_pass http://web-dynamic:80/;	
		}
        location / {
			proxy_pass http://web-static:80;
		}
	}
````
 - Nous n'incluons <b>PLUS</b>  ```/etc/nginx/conf.d/*.conf;``` 

 - Nous déclarons un serveur qui écoute sur le port 80, 
 
 ce dernier 
 -  renvoit tous les urls començant par ```/dyn/``` vers le serveur web dynamique. (en trimant le ```/dyn/``` de l'url, ainsi le serveur dynamic croit que l'url est ```/```).  
Ceci est fait grâce au slash de fin après le port.  
- Renvoit les requêtes ```/``` au serveur statique.

Les noms sont résolus grâce à docker-compose qui met les containers dans le même réseau virtuel et donne un nom à chaque container qui peut être utilisé à la place de son ip.

Docker-compose gère la partie réseau, c'est à dire la résolution de noms à partie du nom du container, mais aussi les ports qui sont ouvert sur un container.  
La commande ```expose``` permet de déclarer un port ouvert sur un container.
Tandis que la commande ```port``` permet de déclarer un port ouvert sur un container ET de faire une correspondance depuis un port de l'hôte. On peut donc voir que seul le service reverse proxy est accessible sur le port 80. Tandis que les deux autres services sont inaccessible depuis l'hôte.
Nous pouvons changer les valeurs de expose dans le dockerfile afin de changer les ports ouverts, mais si on fait ceci, il faudrait aussi changer les valeurs de port dans le fichier nginx.conf dans la directive proxy_pass.
Le dockerfile de l'étape 3 copie juste la configuration du reverse proxy dans le container nginx.

Docker-compose nous permet de ne pas avoir une configuration statique, mais une configuration dynamique. En effet les adresses ip des containers ne sont pas utilisées dans la configuration, nous utilisons les noms des services afin que docker-compose fasse la résolution DNS. Ainsi nous n'avons pas d'ordre d'allumage des services. Cette solution est plus robuste et plus facile à mettre en place, puisque il n'y a pas besoins de faire une configuration après avoir allumé les containers.

# Partie 4, Requêtes Ajax avec JQery
Ici nous avons modifiés les fichiers de l'étape 1.
Plus spécifiquement nous avons 
- Ajouté un ```id="dyntext"``` à une balise afin de pouvoir la chercher et changer avec du javascript.
- Ajouté le fichier dyn.js dans /step1/landing-page/js/
- Inclus le fichier dyn.js (et jquery) dans le fichier index.html.

le Docker-compose de l'étape 3 est utilisable pour voir le résultat.
il suffit d'aller dans le dossier step3 et de lancer  
```sudo docker-compose up -d --build```  
L'option  --build assure que le container est bien construit avec les derniers changement.

En théorie si le Same-Origin-Policy est strict, ça peut poser des problèmes mais puisqu'ici les scripts sont embarqués ça joue. La documentation de mozilla le confirme   

    "Voici quelques exemples de ressources qui peuvent être embarqués malgré leur origine incompatible avec la same-origin policy :

    JavaScript avec <script src="..."></script>. Les messages d'erreur de syntaxe ne sont disponibles que pour les scripts ayant la même origine.
    ...
[source](https://developer.mozilla.org/fr/docs/Web/Security/Same-origin_policy)
# Partie 5, Dynamic reverse proxy configuration

Dans cette partie, nous avons décidé de faire un reverse proxy dynamique à l'aide de [Traefik](https://traefik.io/), il prend les informations directements depuis docker ce qui le rend très puissant puisque il detecte la création de nouveau containers et adapte sa configuration en conséquence.

La configuration de ce dernier se fait directement dans le fichier docker-compose.yml à l'aide de labels.

Ainsi pour que le service web static réponde sur l'url ```/``` nous lui ajoutons lors, de sa description dans le fichier docker-compose.yml, les lignes suivante : 
```yaml
labels:
    - "traefik.http.routers.web-static.rule=PathPrefix(`/`)"
```

Ce que nous faisons ici, c'est créer une route, appelée `web-static` qui est attachée au routeur http qui répond à la règle ```PathPrefix(`/`)``` ainsi les requêtes commençant par ```/``` sont redirigées vers le service web-static.


##  Configuration de web-dynamic: 
`````yaml
labels:
      #create http router web-dynamic with a rule of PathPrefix(`/dyn/`)
      - "traefik.http.routers.web-dynamic.rule=PathPrefix(`/dyn/`)"
      #create the middleware stripPrefix that strips /dyn/ from the requested path
      - "traefik.http.middlewares.stripPrefix.stripprefix.prefixes=/dyn/"
      #register the middleware stripPrefix to the rout web-dynamic
      - "traefik.http.routers.web-dynamic.middlewares=stripPrefix"
`````
Nous faisons ici la mêmes choses : la route ```web-dynamic``` qui répond à la règle ```PathPrefix(`/dyn/`)``` est attachée au routeur http  ainsi les requêtes commençant par ```/dyn/``` sont redirigées vers le service web-dynamic. Cependant, nous créons aussi un middleware qui s'appelle `stripPrefix`, ce dernier enlève le préfixe ```/dyn/``` de l'url. 
Nous attachons maintenant ce middleware à la route web-dynamic.

Si vous souhaitez changer les urls sur lesquels les containers écoutent, il suffit de changer les labels correspondants.

## Configuration de Traefik
Nous utilisons la commande  

```--api.insecure=true --providers.docker```  
Afin d'avoir l'interface web même si nous n'avons pas de certificat SSL, de plus nous annoncons que le fournisseur d'information pour la configuration de Traefik est Docker.

Traefik à besoins d'avoir accès à Docker pour pouvoir déterminer les services qui sont disponibles, ainsi nous devons lui donner un accès au fichier `/var/run/docker.sock` qui est un socket qui permet de communiquer avec le service Docker.
Nous exposons deux ports à l'hôte :
```yaml
ports:
    - "80:80"
    - "8080:8080"
```
80 pour le reverse proxy http et 8080 pour l'interface web de Traefik.


# Additional steps for extra points

## Load balancing multiple server nodes:
Il suffit d'augmenter le nombre de noeuds grâce à docker-compose, traefik se reconfigure automatiquement pour gérer les nouveaux noeuds.
par exemple: 
`````bash
sudo docker-compose up -d --build
sudo docker-compose up -d --scale web-dynamic=3
`````
Grâce à la commande précédente, nous avons créé 3 noeuds web-dynamic.
mainteant en étant sur la page localhost:80, des requêtes ajax sont faites vers les 3 noeuds web-dynamic.
vous pouvez l'observer de la manière suivante : 
```bash 
sudo docker-compose logs -f web-dynamic
```
Vous devriez avoir quelque chose de similaire à :
```text
...
web-dynamic_1    | [GIN] 2022/01/04 - 22:32:57 | 200 |     113.541µs |        10.0.2.2 | GET      "/"
web-dynamic_2    | [GIN] 2022/01/04 - 22:33:02 | 200 |      39.485µs |        10.0.2.2 | GET      "/"
web-dynamic_3    | [GIN] 2022/01/04 - 22:33:07 | 200 |      56.092µs |        10.0.2.2 | GET      "/"
web-dynamic_1    | [GIN] 2022/01/04 - 22:33:12 | 200 |      63.194µs |        10.0.2.2 | GET      "/"
web-dynamic_2    | [GIN] 2022/01/04 - 22:33:17 | 200 |      25.206µs |        10.0.2.2 | GET      "/"
web-dynamic_3    | [GIN] 2022/01/04 - 22:33:22 | 200 |       33.92µs |        10.0.2.2 | GET      "/"
web-dynamic_1    | [GIN] 2022/01/04 - 22:33:27 | 200 |      53.264µs |        10.0.2.2 | GET      "/"
web-dynamic_2    | [GIN] 2022/01/04 - 22:33:32 | 200 |     357.719µs |        10.0.2.2 | GET      "/"
web-dynamic_3    | [GIN] 2022/01/04 - 22:33:37 | 200 |      32.509µs |        10.0.2.2 | GET      "/"
web-dynamic_1    | [GIN] 2022/01/04 - 22:33:42 | 200 |     138.415µs |        10.0.2.2 | GET      "/"
```
On peut aussi observer en utilisant l'interface de Traefik que nous avons bien 3 serveurs pour le service web-dynamic.

## Load balancing: round robin vs sticky sessions
----
Il suffit d'ajouter le label suivant au service statique: 
`````yaml
      - "traefik.http.services.web-static.loadbalancer.sticky.cookie.name=session"
`````
Ainsi nous avons un cookie ayant comme nom session qui est utilié comme sticky cookie afin de déterminer le serveur qui va être utilisé pour la requête.
Nous pouvons regarder les cookies et observer que le cookie session est bien présent sur les requêtes et contient l'id d'un container.  

Pour vérifier que tout fonctionne, nous allons démarrer les 3 services, et augmenter les services `web-static` et `web-dynamic` à 3 noeuds, ensuite nous allons suivre (`logs -f ...`) les logs. Puis dans une premier navigateur nous allons aller sur localhost:80 nous allons voir que les reqêtes dynamiques sont bien réparties sur les 3 noeuds dynamiques, nous pouvons rafraichir la page et voir dans les logs que nous utilisons toujours le même serveur statique. Dans une seconde fenêtre en navigation privée, nous allons aussi nous connecter sur localhost:80. Comme la fenêtre de navigation privée n'envoit pas les cookies, nous devrions utiliser un autre serveur statique. Essayons : 
```bash
sudo docker-compose up -d --scale web-dynamic=3 --scale web-static=3 --build
#Creating network "step5_default" with the default driver
#Creating step5_web-dynamic_1   ... done
#Creating step5_web-dynamic_2   ... done
#Creating step5_web-dynamic_3   ... done
#Creating step5_web-static_1    ... done
#Creating step5_web-static_2    ... done
#Creating step5_web-static_3    ... done
#Creating step5_reverse-proxy_1 ... done

sudo docker-compose logs -f web-dynamic web-static
#initialisation des containers omis par brièveté
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:53:02 +0000] "GET / HTTP/1.1" 304 -
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:53:02 +0000] "GET /js/dyn.js HTTP/1.1" 304 -
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:53:02 +0000] "GET /js/scripts.js HTTP/1.1" 304 -
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:53:02 +0000] "GET /css/styles.css HTTP/1.1" 200 212702
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:53:02 +0000] "GET /assets/img/portfolio-1.jpg HTTP/1.1" 200 254993
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:53:02 +0000] "GET /assets/img/portfolio-3.jpg HTTP/1.1" 304 -
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:53:02 +0000] "GET /assets/img/portfolio-4.jpg HTTP/1.1" 200 240136
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:53:02 +0000] "GET /assets/img/portfolio-2.jpg HTTP/1.1" 200 247590
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:53:02 +0000] "GET /assets/img/bg-masthead.jpg HTTP/1.1" 304 -
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:53:02 +0000] "GET /assets/img/bg-callout.jpg HTTP/1.1" 304 -
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:53:02 +0000] "GET /assets/favicon.ico HTTP/1.1" 200 23462
web-dynamic_1    | [GIN] 2022/01/04 - 22:53:07 | 200 |     210.701µs |        10.0.2.2 | GET      "/"
web-dynamic_3    | [GIN] 2022/01/04 - 22:53:12 | 200 |     147.062µs |        10.0.2.2 | GET      "/"
web-dynamic_2    | [GIN] 2022/01/04 - 22:53:17 | 200 |     117.822µs |        10.0.2.2 | GET      "/"
web-dynamic_1    | [GIN] 2022/01/04 - 22:53:22 | 200 |      43.315µs |        10.0.2.2 | GET      "/"
web-dynamic_3    | [GIN] 2022/01/04 - 22:53:27 | 200 |      99.671µs |        10.0.2.2 | GET      "/"
web-dynamic_2    | [GIN] 2022/01/04 - 22:53:32 | 200 |      42.766µs |        10.0.2.2 | GET      "/"
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:53:33 +0000] "GET / HTTP/1.1" 304 -
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:53:33 +0000] "GET /js/dyn.js HTTP/1.1" 304 -
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:53:33 +0000] "GET /js/scripts.js HTTP/1.1" 304 -
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:53:33 +0000] "GET /assets/img/portfolio-1.jpg HTTP/1.1" 304 -
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:53:33 +0000] "GET /assets/img/portfolio-3.jpg HTTP/1.1" 304 -
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:53:33 +0000] "GET /assets/img/portfolio-2.jpg HTTP/1.1" 304 -
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:53:33 +0000] "GET /assets/img/portfolio-4.jpg HTTP/1.1" 304 -
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:53:33 +0000] "GET /assets/favicon.ico HTTP/1.1" 200 23462
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:53:35 +0000] "GET / HTTP/1.1" 200 11434
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:53:35 +0000] "GET /js/dyn.js HTTP/1.1" 304 -
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:53:35 +0000] "GET /js/scripts.js HTTP/1.1" 304 -
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:53:35 +0000] "GET /assets/img/portfolio-2.jpg HTTP/1.1" 304 -
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:53:35 +0000] "GET /assets/img/portfolio-1.jpg HTTP/1.1" 304 -
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:53:35 +0000] "GET /assets/img/portfolio-3.jpg HTTP/1.1" 304 -
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:53:35 +0000] "GET /assets/img/portfolio-4.jpg HTTP/1.1" 304 -
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:53:35 +0000] "GET /assets/favicon.ico HTTP/1.1" 200 23462
web-dynamic_1    | [GIN] 2022/01/04 - 22:53:40 | 200 |      38.794µs |        10.0.2.2 | GET      "/"
web-dynamic_3    | [GIN] 2022/01/04 - 22:53:45 | 200 |      49.025µs |        10.0.2.2 | GET      "/"
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:53:48 +0000] "GET / HTTP/1.1" 304 -
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:53:48 +0000] "GET /js/dyn.js HTTP/1.1" 304 -
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:53:48 +0000] "GET /js/scripts.js HTTP/1.1" 304 -
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:53:48 +0000] "GET /assets/img/portfolio-1.jpg HTTP/1.1" 200 254993
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:53:48 +0000] "GET /assets/img/portfolio-3.jpg HTTP/1.1" 304 -
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:53:48 +0000] "GET /assets/img/portfolio-4.jpg HTTP/1.1" 200 240136
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:53:48 +0000] "GET /assets/img/portfolio-2.jpg HTTP/1.1" 200 247590
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:53:48 +0000] "GET /assets/favicon.ico HTTP/1.1" 200 23462
web-dynamic_2    | [GIN] 2022/01/04 - 22:53:53 | 200 |      75.648µs |        10.0.2.2 | GET      "/"
web-dynamic_1    | [GIN] 2022/01/04 - 22:53:58 | 200 |       19.53µs |        10.0.2.2 | GET      "/"
web-static_3     | 172.24.0.7 - - [04/Jan/2022:22:54:01 +0000] "GET / HTTP/1.1" 200 11434
web-static_3     | 172.24.0.7 - - [04/Jan/2022:22:54:01 +0000] "GET /js/dyn.js HTTP/1.1" 200 695
web-static_3     | 172.24.0.7 - - [04/Jan/2022:22:54:01 +0000] "GET /css/styles.css HTTP/1.1" 200 212702
web-static_3     | 172.24.0.7 - - [04/Jan/2022:22:54:01 +0000] "GET /js/scripts.js HTTP/1.1" 200 2747
web-static_3     | 172.24.0.7 - - [04/Jan/2022:22:54:01 +0000] "GET /assets/img/portfolio-1.jpg HTTP/1.1" 200 254993
web-static_3     | 172.24.0.7 - - [04/Jan/2022:22:54:01 +0000] "GET /assets/img/portfolio-2.jpg HTTP/1.1" 200 247590
web-static_3     | 172.24.0.7 - - [04/Jan/2022:22:54:01 +0000] "GET /assets/img/portfolio-3.jpg HTTP/1.1" 200 299484
web-static_3     | 172.24.0.7 - - [04/Jan/2022:22:54:01 +0000] "GET /assets/img/portfolio-4.jpg HTTP/1.1" 200 240136
web-static_3     | 172.24.0.7 - - [04/Jan/2022:22:54:01 +0000] "GET /assets/img/bg-callout.jpg HTTP/1.1" 200 1829666
web-static_3     | 172.24.0.7 - - [04/Jan/2022:22:54:01 +0000] "GET /assets/img/bg-masthead.jpg HTTP/1.1" 200 1687843
web-static_3     | 172.24.0.7 - - [04/Jan/2022:22:54:01 +0000] "GET /assets/favicon.ico HTTP/1.1" 200 23462
web-dynamic_3    | [GIN] 2022/01/04 - 22:54:03 | 200 |      16.009µs |        10.0.2.2 | GET      "/"
web-static_3     | 172.24.0.7 - - [04/Jan/2022:22:54:04 +0000] "GET / HTTP/1.1" 304 -
web-static_3     | 172.24.0.7 - - [04/Jan/2022:22:54:04 +0000] "GET /js/dyn.js HTTP/1.1" 304 -
web-static_3     | 172.24.0.7 - - [04/Jan/2022:22:54:04 +0000] "GET /js/scripts.js HTTP/1.1" 304 -
web-static_3     | 172.24.0.7 - - [04/Jan/2022:22:54:04 +0000] "GET /assets/img/portfolio-2.jpg HTTP/1.1" 304 -
web-static_3     | 172.24.0.7 - - [04/Jan/2022:22:54:04 +0000] "GET /assets/img/portfolio-1.jpg HTTP/1.1" 304 -
web-static_3     | 172.24.0.7 - - [04/Jan/2022:22:54:04 +0000] "GET /assets/img/portfolio-3.jpg HTTP/1.1" 304 -
web-static_3     | 172.24.0.7 - - [04/Jan/2022:22:54:04 +0000] "GET /assets/img/portfolio-4.jpg HTTP/1.1" 304 -
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:54:05 +0000] "GET / HTTP/1.1" 200 11434
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:54:05 +0000] "GET /js/dyn.js HTTP/1.1" 304 -
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:54:05 +0000] "GET /js/scripts.js HTTP/1.1" 304 -
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:54:05 +0000] "GET /assets/img/portfolio-2.jpg HTTP/1.1" 200 247590
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:54:05 +0000] "GET /assets/img/portfolio-1.jpg HTTP/1.1" 200 254993
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:54:05 +0000] "GET /assets/img/portfolio-4.jpg HTTP/1.1" 200 240136
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:54:05 +0000] "GET /assets/img/portfolio-3.jpg HTTP/1.1" 304 -
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:54:06 +0000] "GET /assets/favicon.ico HTTP/1.1" 200 23462
web-static_3     | 172.24.0.7 - - [04/Jan/2022:22:54:07 +0000] "GET / HTTP/1.1" 304 -
web-static_3     | 172.24.0.7 - - [04/Jan/2022:22:54:07 +0000] "GET /js/dyn.js HTTP/1.1" 304 -
web-static_3     | 172.24.0.7 - - [04/Jan/2022:22:54:07 +0000] "GET /js/scripts.js HTTP/1.1" 304 -
web-static_3     | 172.24.0.7 - - [04/Jan/2022:22:54:07 +0000] "GET /assets/img/portfolio-1.jpg HTTP/1.1" 304 -
web-static_3     | 172.24.0.7 - - [04/Jan/2022:22:54:07 +0000] "GET /assets/img/portfolio-2.jpg HTTP/1.1" 304 -
web-static_3     | 172.24.0.7 - - [04/Jan/2022:22:54:07 +0000] "GET /assets/img/portfolio-4.jpg HTTP/1.1" 304 -
web-static_3     | 172.24.0.7 - - [04/Jan/2022:22:54:07 +0000] "GET /assets/img/portfolio-3.jpg HTTP/1.1" 304 -
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:54:08 +0000] "GET / HTTP/1.1" 200 11434
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:54:08 +0000] "GET /js/dyn.js HTTP/1.1" 200 695
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:54:08 +0000] "GET /js/scripts.js HTTP/1.1" 304 -
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:54:08 +0000] "GET /assets/img/portfolio-1.jpg HTTP/1.1" 200 254993
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:54:08 +0000] "GET /assets/img/portfolio-2.jpg HTTP/1.1" 200 247590
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:54:08 +0000] "GET /assets/img/portfolio-3.jpg HTTP/1.1" 304 -
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:54:08 +0000] "GET /assets/img/portfolio-4.jpg HTTP/1.1" 200 240136
web-static_2     | 172.24.0.7 - - [04/Jan/2022:22:54:08 +0000] "GET /assets/favicon.ico HTTP/1.1" 200 23462
web-dynamic_2    | [GIN] 2022/01/04 - 22:54:12 | 200 |      16.277µs |        10.0.2.2 | GET      "/"
web-dynamic_1    | [GIN] 2022/01/04 - 22:54:13 | 200 |      18.141µs |        10.0.2.2 | GET      "/"
web-dynamic_3    | [GIN] 2022/01/04 - 22:54:17 | 200 |      17.015µs |        10.0.2.2 | GET      "/"
web-dynamic_2    | [GIN] 2022/01/04 - 22:54:18 | 200 |      14.382µs |        10.0.2.2 | GET      "/"
web-dynamic_1    | [GIN] 2022/01/04 - 22:54:22 | 200 |     225.005µs |        10.0.2.2 | GET      "/"
web-dynamic_3    | [GIN] 2022/01/04 - 22:54:23 | 200 |      35.284µs |        10.0.2.2 | GET      "/"
web-dynamic_2    | [GIN] 2022/01/04 - 22:54:27 | 200 |      54.354µs |        10.0.2.2 | GET      "/"
web-dynamic_1    | [GIN] 2022/01/04 - 22:54:28 | 200 |      57.051µs |        10.0.2.2 | GET      "/"
web-dynamic_3    | [GIN] 2022/01/04 - 22:54:32 | 200 |      49.781µs |        10.0.2.2 | GET      "/"
web-dynamic_2    | [GIN] 2022/01/04 - 22:54:33 | 200 |      36.663µs |        10.0.2.2 | GET      "/"
web-dynamic_1    | [GIN] 2022/01/04 - 22:54:37 | 200 |       70.62µs |        10.0.2.2 | GET      "/"
web-dynamic_3    | [GIN] 2022/01/04 - 22:54:38 | 200 |      50.903µs |        10.0.2.2 | GET      "/"
web-dynamic_2    | [GIN] 2022/01/04 - 22:54:42 | 200 |      40.253µs |        10.0.2.2 | GET      "/"
web-dynamic_1    | [GIN] 2022/01/04 - 22:54:43 | 200 |      39.757µs |        10.0.2.2 | GET      "/"
```
On peut ici observer ce que l'on espérait.
 - Les requêtes pour le service dynamique sont réparties selon un round-robin.
 - Les requêtes pour le service statique sont round-robin mais constante PAR SESSION. 
C'est à dire que pour chaque nouvel utilisateur qui arrive (basé sur le cookie), il ne communiquera qu'avec une même instance du service statique (le même serveur) mais il y a chaque serveur qui recoit la même charge d'utilisateur.

## Dynamique cluster managment
----
Grâce à Traefik c'est fait automatiquement.
Pour s'en convaincre, il suffit de lancer l'interface de Traefik, (port 8080), faire un scale d'un service et de verifier qu'il y a bien plusieurs serveurs pour le service choisis.
````bash
sudo docker-compose up -d --scale web-dynamic=3
````
ou alors de se connecter sur le site, d'attendre que les requêtes dynamiquent commencent, puis de regarder les logs du service web-dynamic avec la commande suivante
````bash
sudo docker-compose logs -f web-dynamic 
````
## Management UI
--- 
Les fichiers relatif au projet sont dans le dossier control-ui.  
La configuration du projet se fait dans le docker-compose du dossier /step5 

Nous avons décidé de faire avec Go et Gin. Nous nous connectons au deamon docker, récupérons les infos et les affichons sur une page web. Le site est derrière le reverse proxy, il est accessible à l'adresse /control/. 

Le choix de l'url est configurable dans Traefik à travers les labels dans le fichier docker-compose.

pour lancer tous les containers avec les dernières améliorations : 
```bash
cd step5
sudo docker-compose up -d 
sudo docker-compose logs -f
```
L'application go écoute par défaut sur le port 80. Le projet à un dockerfile qui permet de containeriser le projet (utilisé par le docker-compose de step5).

Les templates pour construire le site sont dans control-ui/templates.
