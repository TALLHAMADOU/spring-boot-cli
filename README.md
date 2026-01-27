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

```
   ____                  _             ____ _     _
  / ___| _ __   ___  ___| |_ ___ _ __ / ___| |__ (_)_ __   __ _
  \___ \| '_ \ / _ \/ __| __/ _ \ '__| |   | '_ \| | '_ \ / _` |
   ___) | |_) |  __/ (__| ||  __/ |  | |___| | | | | | | | | (_| |
  |____/| .__/ \___|\___|\__\___|_|   \____|_| |_|_|_| |_|\__, |
        |_|                                                |___/

Un petit outil CLI pour générer rapidement des composants Spring Boot (entités JPA, repositories, services, controllers, DTOs, tests) et aider le versioning Maven/Gradle.

Badges
-------

![build](https://img.shields.io/badge/build-local-lightgrey)
![license](https://img.shields.io/badge/license-MIT-blue)

Introduction
------------

Ce CLI est conçu pour être installé globalement (une seule fois) et utilisé depuis n'importe quel dossier projet. Toutes les opérations modifient uniquement le dossier courant. Il facilite la génération de code répétitif et la gestion simple de versions pour les projets Maven/Gradle.

Installation
------------

1) Depuis les sources (machine de développement):

```bash
git clone <votre-fork>
cd Spring-CLi
go build -o spring-cli .
# déplacer le binaire dans votre PATH
sudo mv spring-cli /usr/local/bin/
```

2) Compiler en intégrant une version (recommandé pour releases):

```bash
go build -ldflags "-X 'github.com/hamadoutall/spring-cli/cmd.Version=1.2.3'" -o spring-cli .
```

3) Utiliser `goreleaser` pour publier des releases multiplateformes (le fichier `goreleaser.yml` est fourni et peut être adapté).

Affichage de la version du CLI
----------------------------

Vous pouvez obtenir la version du binaire :

```bash
spring-cli --version    # affiche la version (Cobra)
spring-cli -v           # alias court, affiche la version du CLI
```

Utilisation (rappels rapides)
-----------------------------

Toutes les commandes s'exécutent depuis la racine du projet Spring (ou d'un dossier où vous voulez générer le code).

Exemples rapides :

```bash
# créer un projet Maven minimal
spring-cli install:project maven --name testproj --package com.example.test

# générer une entité User avec deux champs
spring-cli make entity User --fields "name:String,age:int"

# générer repository/service/controller pour User
spring-cli make repository User
spring-cli make service User --entity User
spring-cli make controller User --entity User --crud

# inspecter / modifier la version du projet
spring-cli version            # affiche la version détectée
spring-cli version --bump minor
spring-cli version --set 0.2.0
```

Commandes et flags détaillés
---------------------------

- `install:project <maven|gradle>`
  - Flags: `--name, -n` (nom du projet), `--package, -p` (package de base)
  - Crée un projet minimal (Application.java + build file + wrappers)

- `make entity NAME`
  - Flags: `--fields` (ex: `id:Long,name:String`), `--lombok`, `--auditing`, `--package, -p`
  - Ajoute automatiquement la dépendance JPA si absente

- `make repository NAME`
  - Flags: `--package, -p`

- `make service NAME`
  - Flags: `--entity, -e`, `--package, -p`

- `make controller NAME`
  - Flags: `--crud`, `--entity`, `--package, -p`

- `make dto NAME` et `make request NAME`
  - Flags: `--fields`, `--package, -p`

- `make test <service|controller>`
  - Génère des templates de tests JUnit sous `src/test/java`.

Versioning (commande `version`)
--------------------------------

La commande `version` détecte la version dans `pom.xml` (Maven) ou `build.gradle` (Gradle).

Flags principaux:

- `--bump patch|minor|major` : incrémente la version détectée.
- `--set x.y.z` : fixe la version explicitement.
- `--auto` : récupère le dernier tag Git (ex: `v1.2.3`), l'incrémente (patch) et propose la nouvelle version.
- `--commit` : `git add` + `git commit` des fichiers `pom.xml`/`build.gradle` modifiés.
- `--tag` : crée un tag `vX.Y.Z`.
- `--push` : pousse les commits / tags vers le remote (si configuré).

Sécurité et comportement Git
---------------------------

- `--commit`, `--tag` et `--push` sont des actions explicites : elles ne s'exécutent que si vous demandez.
- Si le dossier n'est pas un dépôt Git, la commande `--commit` affichera une erreur Git claire plutôt que d'initialiser automatiquement un repo.
- `--push` vérifiera l'existence d'un remote ; si aucun remote n'est configuré, elle affichera un message et n'essaiera pas de pousser.

Détection du package de base
---------------------------

Le CLI tente de déterminer le package de base du projet via :

1. Le flag `--package` passé à la commande actuelle
2. Le flag `--package` passé lors d'un `install:project`
3. La lecture de `pom.xml` (`groupId`) ou `build.gradle` (`group`) et `settings.gradle`

Cette logique permet d'éviter d'écrire dans `com.example` par défaut quand votre projet utilise un package personnalisé.

Personnalisation visuelle (logo)
-------------------------------

Le README inclut un petit ASCII-art ci-dessus pour l'identité. Si vous préférez utiliser le logo officiel Spring, placez le fichier `logo.png` dans `docs/` et référencez-le dans vos assets/documentation (respectez la licence et la marque Spring si vous utilisez le logo officiel).

Conseils pour les releases
-------------------------

- Préparez votre `goreleaser.yml` (un fichier d'exemple est fourni). Exécutez `goreleaser release --rm-dist` depuis un tag Git.
- Construisez localement pour tester :

```bash
go build -ldflags "-X 'github.com/hamadoutall/spring-cli/cmd.Version=0.0.0'" -o spring-cli .
```

Contribuer
----------

- Ouvrez une issue pour proposer des améliorations ou signaler des bugs.
- Forkez, changez et proposez une PR. Respectez le style du dépôt et ajoutez des tests pour les nouvelles fonctionnalités.

Foire aux questions (FAQ)
------------------------

Q: Le CLI modifie-t-il des fichiers en dehors du dossier courant ?
A: Non — tout est local au dossier courant.

Q: Puis-je utiliser `spring-cli` en CI ?
A: Oui. Pour CI, construisez le binaire dans votre pipeline ou utilisez une release GitHub fournie par `goreleaser`.

Licence
-------

Voir le fichier `LICENSE` pour les détails.

Fichiers utiles
---------------

- [cmd/version.go](cmd/version.go) : implémentation de la commande `version`.
- [goreleaser.yml](goreleaser.yml) : configuration de release (modifiez owner/repo avant usage).

***

Merci d'utiliser `spring-cli` — dites-moi si vous voulez que j'ajoute des captures d'écran, des exemples de templates Java plus riches, ou un site de documentation généré.

