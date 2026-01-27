# spring-cli

Un petit CLI pour générer rapidement des composants Spring Boot (entities, services, repositories, etc.).

Installation (binaire global)

- Avec Go (machine de développement) :

```bash
go build -o spring-cli .
# puis déplacer le binaire dans un dossier du PATH, par ex:
sudo mv spring-cli /usr/local/bin/
```

- Utiliser les releases (Linux/macOS/Windows) :
  - Télécharger l'archive depuis les releases GitHub et placer le binaire dans un dossier du `PATH`.

Utilisation

- Générer une entité dans le dossier courant (ne modifie que le projet courant) :

```bash
cd mon-projet-spring
spring-cli make entity User --fields "name:String,age:int"
```

- Générer un service lié à une entité :

***
# spring-cli
 # spring-cli

 ![license](https://img.shields.io/badge/license-MIT-blue)

    ____                  _             ____ _     _
   / ___| _ __   ___  ___| |_ ___ _ __ / ___| |__ (_)_ __   __ _
   \___ \| '_ \ / _ \/ __| __/ _ \ '__| |   | '_ \| | '_ \ / _` |
    ___) | |_) |  __/ (__| ||  __/ |  | |___| | | | | | | | | | (_| |
   |____/| .__/ \___|\___|\__\___|_|   \____|_| |_|_|_| |_|\__, |
         |_|                                                |___/

 Un outil CLI pour générer rapidement des composants Spring Boot (entités JPA, repositories, services, controllers, DTOs, tests) et aider le versioning pour projets Maven/Gradle.

 ## Table des matières
 - Introduction
 - Installation
 - Utilisation rapide
 - Commandes et flags
 - Versioning
 - Exemples : création de projet (Maven/Gradle)
 - Contribuer
 - Licence

 ## Introduction

 Ce CLI s'installe globalement et s'utilise depuis n'importe quel dossier projet. Toutes les opérations modifient uniquement le dossier courant et visent à accélérer la génération de code répétitif.

 ## Installation

 1) Depuis les sources (machine de développement)

 ```bash
 git clone <votre-fork>
 cd Spring-CLi
 go build -o spring-cli .
 # déplacer le binaire dans votre PATH
 sudo mv spring-cli /usr/local/bin/
 ```

 2) Compiler en intégrant une version (pour releases)

 ```bash
 go build -ldflags "-X 'github.com/hamadoutall/spring-cli/cmd.Version=1.2.3'" -o spring-cli .
 ```

 3) Utiliser les releases GitHub

 Téléchargez l'archive depuis les releases et placez le binaire dans un dossier présent dans votre `PATH`.

 4) Publier avec `goreleaser`

 Le dépôt contient un fichier `goreleaser.yml` prêt à être adapté. Exécutez `goreleaser release --rm-dist` depuis un tag Git.

 ## Utilisation rapide

 Exemples courants (exécutés depuis la racine du projet Spring) :

 ```bash
 # créer un projet Maven minimal
 spring-cli install:project maven --name testproj --package com.example.test

 # créer un projet Gradle minimal
 spring-cli install:project gradle --name testproj --package com.example.test

 # générer une entité User avec deux champs
 spring-cli make entity User --fields "name:String,age:int"

 # générer repository/service/controller pour User
 spring-cli make repository User
 spring-cli make service User --entity User
 spring-cli make controller User --entity User --crud

 # gestion de version
 spring-cli version            # affiche la version détectée dans pom.xml/build.gradle
 spring-cli version --bump minor
 spring-cli version --set 0.2.0
 ```

 ## Commandes et flags (résumé)

 - `install:project <maven|gradle>`
   - Flags: `--name, -n` (nom du projet), `--package, -p` (package de base)
   - Crée un projet minimal (Application.java/Kotlin + fichier de build + wrappers)

 - `make entity NAME`
   - Flags: `--fields` (ex: `id:Long,name:String`), `--lombok`, `--auditing`, `--package, -p`

 - `make repository NAME` — Flags: `--package, -p`
 - `make service NAME` — Flags: `--entity, -e`, `--package, -p`
 - `make controller NAME` — Flags: `--crud`, `--entity`, `--package, -p`
 - `make dto NAME` / `make request NAME` — Flags: `--fields`, `--package, -p`
 - `make test <service|controller>` — Génère des templates JUnit sous `src/test/java`.

 ## Versioning (`version`)

 La sous-commande `version` détecte la version dans `pom.xml` (Maven) ou `build.gradle` (Gradle).

 Flags utiles:
 - `--bump patch|minor|major` : incrémente la version
 - `--set x.y.z` : fixe la version
 - `--auto` : récupère le dernier tag Git et propose une incrémentation
 - `--commit`, `--tag`, `--push` : actions Git explicites (exécutées seulement si demandées)

 ## Exemples : création d'un projet (Maven / Gradle)

 ### Maven (création rapide)

 ```bash
 spring-cli install:project maven --name demo-maven --package com.example.demo
 ```

 Résultat attendu (exemples de fichiers créés) :
 - `pom.xml`
 - `src/main/java/com/example/demo/DemoApplication.java`
 - `src/main/resources/application.properties`
 - wrappers `mvnw`, `mvnw.cmd` et `.mvn/`

 ### Gradle (création rapide)

 ```bash
 spring-cli install:project gradle --name demo-gradle --package com.example.demo
 ```

 Résultat attendu (exemples de fichiers créés) :
 - `build.gradle` ou `build.gradle.kts`
 - `settings.gradle` ou `settings.gradle.kts`
 - `src/main/java/com/example/demo/DemoApplication.java`
 - wrappers `gradlew`, `gradlew.bat` et `gradle/`

 ## Sécurité et comportement Git

 Les actions `--commit`, `--tag` et `--push` sont exécutées uniquement si vous les spécifiez. Si le dossier n'est pas un dépôt Git, la commande affichera une erreur claire.

 ## Contribuer

 Ouvrez une issue pour proposer des améliorations ou signaler un bug. Forkez, modifiez et proposez une PR en respectant le style du dépôt et en ajoutant des tests si nécessaire.

 ## Fichiers utiles

 - [cmd/version.go](cmd/version.go) : implémentation de la commande `version`.
 - [goreleaser.yml](goreleaser.yml) : configuration de release (à adapter owner/repo).

 ## Licence

 Voir le fichier `LICENSE` pour les détails.

 Merci d'utiliser `spring-cli` — dites-moi si vous voulez que j'ajoute des captures d'écran, des exemples de templates Java plus riches, ou de la documentation générée.
go build -ldflags "-X 'github.com/hamadoutall/spring-cli/cmd.Version=0.0.0'" -o spring-cli .

