# spring-cli

Un petit CLI pour générer rapidement des composants Spring Boot (entities, services, repositories, etc.).

Installation (binaire global)
# spring-cli

[![license](https://img.shields.io/badge/license-MIT-blue)](LICENSE) [![go](https://img.shields.io/badge/go-1.25+-00ADD8)

### Résumé

`spring-cli` est un petit outil en Go pour générer rapidement des composants Spring Boot (entités JPA, repositories, services, controllers, DTOs, tests) et pour gérer simplement la version d'un projet Maven/Gradle.

---

## Table des matières

- [Installation](#installation)
- [Utilisation rapide](#utilisation-rapide)
- [Commandes principales](#commandes-principales)
- [Versioning](#versioning)
- [Exemples : création de projet (Maven / Gradle)](#exemples--cr%C3%A9ation-de-projet-maven--gradle)
- [Contribuer](#contribuer)
- [Licence](#licence)

---

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

Le dépôt contient `goreleaser.yml`. Pour publier :

```bash
goreleaser release --rm-dist
```

---

## Utilisation rapide

Exemples depuis la racine d'un projet Spring :

```bash
# créer un projet Maven minimal
spring-cli install:project maven --name testproj --package com.example.test

# créer un projet Gradle minimal
spring-cli install:project gradle --name testproj --package com.example.test

# générer une entité User
spring-cli make entity User --fields "name:String,age:int"

# générer repository/service/controller pour User
spring-cli make repository User
spring-cli make service User --entity User
spring-cli make controller User --entity User --crud

# afficher / modifier la version
spring-cli version
spring-cli version --bump minor
spring-cli version --set 0.2.0
```

---

## Commandes principales

- `install:project <maven|gradle>`
  - Flags: `--name, -n` (nom), `--package, -p` (package de base)
- `make entity NAME` — Flags: `--fields`, `--lombok`, `--auditing`, `--package`
- `make repository NAME` — Flags: `--package`
- `make service NAME` — Flags: `--entity`, `--package`
- `make controller NAME` — Flags: `--crud`, `--entity`, `--package`
- `make dto NAME` / `make request NAME` — Flags: `--fields`, `--package`
- `make test <service|controller>` — Génère des templates de tests JUnit sous `src/test/java`

---

## Versioning (`version`)

Détecte la version dans `pom.xml` (Maven) ou `build.gradle` (Gradle).

Flags utiles :

- `--bump patch|minor|major` — incrémente la version
- `--set x.y.z` — fixe la version
- `--auto` — récupère le dernier tag Git et propose une incrémentation
- `--commit`, `--tag`, `--push` — actions Git explicites (s'exécutent seulement si demandées)

---

## Exemples : création de projet (Maven / Gradle)

### Maven

```bash
spring-cli install:project maven --name demo-maven --package com.example.demo
```

Fichiers créés (exemple) :

- `pom.xml`
- `src/main/java/com/example/demo/DemoApplication.java`
- `src/main/resources/application.properties`
- wrappers: `mvnw`, `mvnw.cmd`, `.mvn/`

### Gradle

```bash
spring-cli install:project gradle --name demo-gradle --package com.example.demo
```

Fichiers créés (exemple) :

- `build.gradle` or `build.gradle.kts`
- `settings.gradle` or `settings.gradle.kts`
- `src/main/java/com/example/demo/DemoApplication.java`
- wrappers: `gradlew`, `gradlew.bat`, `gradle/`

---

## Contribuer

- Ouvrez une issue pour une suggestion ou un bug.
- Forkez, modifiez et envoyez une PR en respectant le style du dépôt et en ajoutant des tests si nécessaire.

## Fichiers utiles

- [cmd/version.go](cmd/version.go) — implémentation de la commande `version`.
- [goreleaser.yml](goreleaser.yml) — configuration de release (à adapter).

## Licence

Voir le fichier `LICENSE`.

---

Merci d'utiliser `spring-cli`. Dites-moi si vous préférez un style différent (par ex. README en anglais, captures d'écran, ou documentation MkDocs). 
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

 Merci d'utiliser `spring-cli`
