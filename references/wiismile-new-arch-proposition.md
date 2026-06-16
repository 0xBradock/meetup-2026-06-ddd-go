TL;DR

“Make the change easy, than make the easy change” - Kent Beck

Une refonte augmente la vitesse du business en permettant une livraison plus rapide et plus efficace des fonctionnalités, tout en réduisant les coûts de maintenance et en améliorant la qualité globale des applications.

-   Strangler Fig - La refonte doit se faire par étapes (https://martinfowler.com/bliki/StranglerFigApplication.html ). Nous pouvons commencer par le module le plus simple pour comprendre les complexités inhérentes et cachées; ensuite, on peut passer à d’autres modules. Cela permettra également de créer des templates et modules réutilisables, ce qui se fait facilement dans une architecture découpée en domaine (DDD).
-   Tests (dev + QA) - Sont cruciaux et doivent être utilisés comme guide lors de la refonte. Si le code existant ne possède pas de tests, nous pouvons utiliser les tests qui seront créés par le nouveau code.
-   DDD - Est une pratique qui regroupe des activités côté business (tactique) et tech (stratégique) bien documenté et permet d'aligner au mieux le business avec la tech.
-   SonarCloud - Toute nouvelle refonte doit être accompagné d’une stratégie de test qui peut être utilisé pour valider contre le système ancien (Clean as You Code - SonarCloud)
-   Technos - Go et PHP peuvent co-exister dans nos techos backend. Une alternative basée sur JS pour le front (comme Next.js , React Router ) est conseillée.
-   Templates - Nous pouvons commencer par identifier les blocs de code réutilisables, ce qui garantit une livraison plus rapide pour les futures refontes ou refactorings.

Un découpage par domaine présente plusieurs avantages et permet de mieux aligner le business avec le dev

Je me concentre uniquement sur les aspects techniques et les interactions autour du dev technique. Il y a tout un travail personnel et de gestion humaine que je n’aborderai pas ici, et qui reste à traiter en collaboration avec d’autres membres de l’équipe.

1 Contexte & Retours d’Expérience

Réflexion sur des initiatives passées (en interne ou observées ailleurs) pour nourrir la décision.

Rappel de l'architecture proposée.

Ce document présente une revue complète de l'architecture actuelle et une proposition de stratégie de réécriture progressive. En s’appuyant sur une analyse d’alternatives architecturales détaillées (cf. annexe A1), il met en lumière les faiblesses de l’existant et les enseignements tirés d’expériences passées. L’architecture cible est examinée sous l’angle de ses forces, faiblesses, opportunités et menaces. La stratégie recommandée repose sur un découpage en deux étapes : un travail préparatoire sur les tests, les métriques et le déploiement, suivi d’une migration incrémentale module par module. Enfin, des choix technologiques sont proposés pour le frontend et le backend, avec un accent sur la scalabilité, l’éco-conception et l’attractivité pour les développeurs.

1.1 Pour Quoi Faire la Refonte

“Reduce volatility in the marginal cost of features” - J.B. Rainsberg

Au début, la création de code semble simple et directe, mais avec le temps, la complexité s'accumule, rendant la maintenance de plus en plus coûteuse et complexe, révélant ainsi que le véritable coût d'un logiciel réside dans sa maintenance à long terme.

Une décomposition de nos applications par domaine (DDD) me semble le bon compromis entre compléxité et isolation. Cette approche permet de structurer le code en modules distincts, chacun représentant un domaine métier spécifique, ce qui facilite la maintenance, améliore la scalabilité et réduit la dette technique. En adoptant une architecture basée sur des domaines, nous pouvons non seulement simplifier la gestion du code, mais aussi rendre les futures modifications plus faciles et moins risquées, tout en améliorant les performances et la sécurité de l'application.

En conclusion, l'objectif final d'une refonte est d'augmenter la vitesse du business en permettant une livraison plus rapide et plus efficace des fonctionnalités, tout en réduisant les coûts de maintenance et en améliorant la qualité globale de l'application

1.2 Points d’échec / Enseignements

Les points qui peuvent faire échouer le projet

-   Manque de tests:
    Netscape, JoelOnSoftware 2 - 2002 (When they did, finally, release their software, it doesn’t seem like they did very much testing.)
-   Envie de réécriture par manque de comprehension:
    Michael Feathers (Usually, this is because the reason that people want to rewrite code is because they don't understand it. Yet, rewriting code often requires us to understand it well enough to proceed with the rewrite, especially if there are existing customers who depend on all of the nuances of behavior that the system has consistently exhibited.)
-   Big Bang:
    Netscape, JoelOnSoftware - 2000 (They did it by making the single worst strategic mistake that any software company can make: They decided to rewrite the code from scratch.)Netscape, JoelOnSoftware 2 - 2002 (Lou Montulli, one of the 5 programming superstars who did the original version of Navigator ... "... it’s one of the major reasons I resigned from Netscape." This one decision cost Netscape 3 years.)Netscape, JoelOnSoftware - 2002 (They shouldn’t have rewritten from scratch. They should have done this all in steps. Big chunky steps, fine, but steps. For example, they could have rebuilt the rendering engine — without touching any of the other stuff — as a first step. Then ship. Was there anything wrong with the networking library? I don’t think there was. Even if there was, OK. So, fix it. One step at a time.)
-   Problème de scope:
    Borland (Borland made the same mistake when they bought Arago and tried to make it into dBase for Windows, a doomed project that took so long that Microsoft Access ate their lunch, then they made it again in rewriting Quattro Pro from scratch and astonishing people with how few features it had),

1.3 Cas de réussite

Les points qui contribuent à la réussite

-   Teste en amont:
    TursoSqlite (doing it with Deterministic Simulation Testing (DST) built-in from the get-go),DaedTech (_ Automated tests are your friend — characterize the system’s current behavior with lots of automated tests and then work on refactoring._)
-   Déployment régulier:
    Adobe, Hardy - 2019 (When our general availability release date came, it was largely a formality. There were no surprises to our users. There were no new features revealed. In fact, there was no code release at all. It was uneventful and low-stress, which is exactly what we desired.)
-   Cohésion d'équipe:
    Adobe, Hardy - 2019 (You must have organizational buy-in from the top to the bottom. Communicate the cost and benefits clearly. Have a plan. Define success and failure. Don’t attempt to hide the difficulty or investment required.)

Comment éviter le Besoin d'une Autre Refonte Future

Refactor continu, attention à l’architecture et adoption d’une stratégie de test robuste.

-   DaedTech (The software is a mess because the group made it a mess, and it’ll only get and stay clean if the group cleans it.)
-   Netscape, JoelOnSoftware - 2000 (Architecture, inneficiency and careless code)

2 Analyse de l'Architecture Proposée

Retour critique sur l’architecture proposée.

Options d'Architecture

Forces

-   Encapsulation des domaines fonctionnels (domain-driven design potentiel)
-   Meilleure scalabilité et résilience
-   Équipes plus autonomes
-   Utilisation de technologies familières à l'équipe et attractives pour les nouveaux recrutements (ex : PHP Symfony, Go et React/Vue)

Faiblesses

-   Complexité accrue (infrastructure, communication inter-service)
-   Manque de stratégie de test initiale
-   Courbe d’apprentissage pour des profils moins expérimentés
-   Risque d’over-engineering dès le début
-   Cohérence technique difficile à maintenir sans normes partagées
-   Fautes classiques des systèmes distribués : le réseau est fiable, latence nulle, etc.
-   Expérience avec une nouvelle organisation d’équipe qui se met en place

Mesures pour mitiger les faiblesses : pair programming, ...

Opportunités

-   Utiliser des techniques pour mieux aligner les experts du domaine (PO) avec les devs
-   Renforcer la culture de l’automatisation et du test
-   Intégration de la culture DevOps et acquisition des compétences d’infrastructure d’application
-   Accélérer les cycles de développement grâce à des livraisons indépendantes
-   Améliorer la fiabilité et la résilience (grâce à l’observabilité & l’isolation)
-   Structurer la dette technique de façon plus maîtrisée
-   Favoriser la montée en compétence de l’équipe

Menaces

-   Réécriture partielle ou totale difficile à tenir en termes de planning
-   Effet tunnel si la transition n’est pas progressive
-   Risque d’avoir une cohabitation instable entre legacy et nouveau code
-   Temps d’onboarding plus long pour de nouveaux développeurs
-   Dépendance excessive à des technologies peu maîtrisées en interne
-   Coûts d’infrastructure non maîtrisés (AWS)

Mesures pour mitiger les menaces : formation AWS, avoir des équipes dédiées pour la réécriture, ...

3 Stratégie de Réécriture Proposée

Plutôt qu’une réécriture "big bang", je recommande une évolution itérative et contrôlée, axée sur des tests.

Un découpage par domaine présente plusieurs avantages et permet de mieux aligner le business avec le dev

Chaque element découpé du serveur doit être independent de ce que lui déclenche et de ce qui est déclenché. Ce-la permet d’avoir des elements modulaires qui peuvent être repris ?????

Le plan proposé se base sur les principes du [Strangler Fig Pattern de Fowler], avec les étapes suivantes :

Travail en amont des modifications :

-   3.1 Métriques de performance
-   3.2 Couverture de tests
-   3.3 Identification d’un module susceptible d'être remplacé
-   3.4 Cohérence des données

Travail de mise à jour (d’un module) :

-   3.5 La mise à jour du système doit être réalisée un module à la fois
-   3.6 Sessions d’EventStorming avec les experts métier, PO et équipe technique
-   3.7 Vérification des métriques
-   3.8 Isoler le module en cours de modification
-   3.9 Déploiement régulier et incrémental

Je donne plus de références, précisions et éléments de contexte sur ces sujets ci-dessous.

3.1 Établissement des métriques de performance, qui justifient le changement

-   Établissement des métriques de performance, justifiant le changement.

Nous pouvons utiliser:

-   Mean Time To Repair
-   Mean Time To Recovery?
-   Mean Time To Ack?
-   Temps de chargement des pages (lighthouse)

???

3.2 Augmentation de la couverture de tests, utilisée pour valider et piloter la réécriture

-   Augmentation de la couverture de tests, utilisée pour valider et piloter la réécriture.

Nous devons utiliser les tests en place pour valider :

-   Possibilité de créer de nouveaux tests pour le code existant, mais ce n’est pas obligatoire pour déclencher la refonte.
    J’essaierai d’aborder ce sujet avec Benjamin pour que nous puissions augmenter la couverture de test des modules XDI
-   Utiliser les tests en place pour le nouveau système
-   Organiser un système permettant de créer rapidement de nouveaux tests (EventStorming + Gherkin)
-   S’appuyer sur des outils de qualité de test comme Sonar et Squash
-   Les modifications sur l’ancien système ne doivent être prises en compte que si elles sont accompagnées de tests.
    Dans ce cas, ces mêmes tests seront utilisés pour valider la nouvelle version

3.3 Identification d’un module susceptible d'être remplacé

3.4 Cohérence des données

-   Détermination du scope des données (context boundary) entre les environnements et comment minimiser les divergences

3.5 La mise à jour du système doit être faite un module à la fois

La refonte de type big bang n’est conseillée par aucun retour d’expérience.
Le démarrage par le module le plus simple et le moins critique permettrait de raffiner l’implémentation et de tester le nouvel écosystème.

3.6 Des sessions d’EventStorming avec les experts du domaine, PO et équipe tech

L’EventStorming (ou charte SICOP) permettrait de mieux aligner les besoins du business avec l’implémentation.
Ces sessions peuvent aussi apporter de la clarté sur des cas qui n’ont pas encore été traités.

3.7 Vérification des métriques

Validation et vérification que les métriques établies ont bien été atteintes.

3.8 Isoler le module en cours de modification

Minimiser le rajout de fonctionnalités dans le module en cours de mise-à jour. Si besoin, ajouter les modifications dans l’ancien système, accompagnées de tests qui seront également utilisés dans le nouveau système.

3.9 Déploiement régulier

Diminuer la boucle de développement et de déploiement.

4 Choix Technologique

4.1 Frontend

Critère

Symfony (Twig + Vue)

Next.js (React)

Découplage du backend

⚠️ Frontend séparé, mais introduit un second langage et une techno supplémentaire côté frontend

✅ Découplage total avec un écosystème unifié

Utilisation de composants dynamiques

⚠️ Possible, mais plus complexe à intégrer proprement

✅ Natif avec React, composants 100 % contrôlés

Intégration d’un Design System

⚠️ Faisable mais demande du sur-mesure (Twig + composants)

✅ Excellent support via Storybook, design tokens, composants

Uniformité UI entre projets

⚠️ Dépend des efforts d’harmonisation entre Twig et Vue

✅ Forte cohérence possible grâce à des packages partagés

Scalabilité

⚠️ Moyenne : rendu serveur, intégration JS partielle

✅ Élevée : SSR/ISR, CDN, découpage logique

Maintenabilité

⚠️ Couplage Vue/Twig parfois fragile, tooling limité

✅ Haute : composants réutilisables, DX moderne

Attractivité pour les devs frontend

❌ Faible : peu de développeurs souhaitent travailler avec Twig

✅ Forte : React est très demandé

Expérience développeur

⚠️ Moins d’outils modernes, intégration JS manuelle

✅ Excellente : tooling moderne, hot reload, TypeScript

SEO

✅ Bon via SSR de Symfony

✅ Excellente via SSR/ISR avec Next.js

💡 Recommandation :

Adoption d’un framework JavaScript/TypeScript avec support pour les composants gérés côté serveur, comme, React Router, Nuxt ou Next.

Utilisation d’un design system pour la gestion d’une bibliothèque de composants pouvant être utilisée à travers toutes les applications frontend (navigateur, mobile, back-office).

Adopter un framework JavaScript moderne permettrait :

-   Une meilleure scalabilité de l’architecture (frontend découplé)
-   Une expérience développeur moderne, attractive pour les profils frontend
-   Une meilleure maintenabilité à long terme grâce à une architecture modulaire
-   D’ouvrir la porte à une organisation plus autonome des équipes frontend/backend
-   Bénéficier d’un écosystème riche avec plusieurs bibliothèques open-source orientées B2B :

-   shadcn (⭐ 85k)
-   daisyUI (⭐ 36k)
-   refine (⭐ 30k)
-   marmelab / react-admin (⭐ 25k)
-   ainsi que des solutions premium comme Tailwind UI Plus (~850 € en licence unique)

L’implémentation d’un design system permettrait :

-   De créer une librairie de composants UI centralisée (via Storybook), accessible à tous les projets
-   D’exposer les composants sous forme de packages NPM privés ou de sous-modules Git
-   D’utiliser des design tokens pour garantir la cohérence des couleurs, typographies, espacements, etc.

4.2 Backend

Critère

Symfony

Go (Golang)

Performance native

⚠️ Interprété, moins performant pour les I/O intensifs

✅ Très performant, compilé, goroutines très efficaces

Scalabilité

⚠️ Nécessite un scaling horizontal plus agressif

✅ Meilleure gestion de la concurrence, faible empreinte mémoire

Éco-conception (consommation serveur)

❌ Charge CPU/mémoire plus élevée pour charges équivalentes

✅ Empreinte écologique réduite grâce à une meilleure efficacité

Coût serveur & infra

❌ Plus de ressources nécessaires ⇒ coût plus élevé

✅ Moins de ressources nécessaires ⇒ réduction des coûts

Expérience développeur

✅ Riche et complexe avec Composer, Symfony CLI, Flex, MakerBundle, etc.

✅ Excellente : outillage natif (format, lint, tests, build), simple à prendre en main

Montée en compétence

✅ Rapide pour développeurs PHP ; écosystème très documenté

✅ L’apprentissage est simple, avec des patterns Go spécifiques mais accessibles

Recrutement

✅ Bon vivier PHP, facile à staffer rapidement

⚠️ Moins de profils disponibles mais très qualifiés et motivés ; montée en compétence facilitée

Modularité du code

⚠️ Symfony encourage de gros packages et une structure MVC classique

✅ Go favorise des services compacts, autonomes, faciles à maintenir

Déploiement

⚠️ Dépendances serveur/web à gérer (FPM, PHP runtime)

✅ Binaire statique, facile à déployer dans des conteneurs, support natif chez AWS

Observabilité / instrumentation

⚠️ Nécessite l’ajout de librairies externes (ex : Monolog, DataDog)

✅ Prometheus, OpenTelemetry : intégration native ou simplifiée

Centre de compétence technique

✅ Facile à construire rapidement avec l’existant

✅ Permet d’attirer des profils "tech passionnés", proches des univers DevOps/SRE

💡 Recommandation :

-   PHP et Go : Possibilité d’une architecture modulaire polyglotte (PHP + Go), évitant une dépendance exclusive à un langage tout en tirant parti des forces de chacun.
-   Go pour l’éco-conception : Une étude universitaire a comparé l’efficacité énergétique de plusieurs langages sur des problèmes standards. Résultat : Go figure parmi les langages les plus sobres, tandis que PHP consomme nettement plus d’énergie et de mémoire. Ce constat souligne que le choix du langage impacte directement l’empreinte carbone d’un service, indépendamment de la qualité du code. Opter pour une stack plus performante n’est donc pas qu’une décision technique — c’est aussi un acte responsable.

-   Intérêt de Go pour les devs PHP : Diversification du profil, élargissement des opportunités de carrière, compréhension d’un autre paradigme de programmation.

4.3 Architecture

A1. Alternative Architectures

Mono x Distributed

-   Monolithic: Layered, Pipeline and microkernel
-   Distributed: Service-Based, event-based, space-based, service-oriented and microservices

Falacies of distributed systems:

-   Network is reliable
-   Latency is zero
-   Bandwidth is infinite
-   Network is secure
-   The topology never changes
-   There is only one admin
-   Transport cost is zero
-   Network is homogenius

Other distributed issues:

-   Logging management
-   Distributed transactions
-   Development environment
-   API contract maintenance and versioning

Concepts:

-   Agility: Ability of a software system to adapt quickly to changes, reducing time-to-market and lowering costs of change
-   Availability: How long the system is up and running
-   Deployability: Ease with which a software system can be released to production, reducing deployment effort, risks, and associated operational costs
-   Reliability: Degree to which a system functions under specified conditions for a specific period of time
-   Scalability: Hability to increase machine size (vertical) and to add more copies of the same machine (horizontal)
-   Testability: Ease with which a system can be tested for correctness, reducing debugging time and lowering development and maintenance costs
-   Maintainability: Effectiveness and effciency which developers can modify the software

Architecture en couches (Layered Architecture)

Architecture structurée en couches techniques empilées (UI, logique métier, données), où chaque couche dépend de la précédente.

L’une des architectures les plus répandues.

Points positifs :

-   Faible coût et simplicité
-   Facile à comprendre, donc adaptée aux petites applications
-   Un bon choix au début, quand la logique métier n’est pas encore bien définie

Points négatifs :

-   Lorsque l’application grandit, la maintenabilité, l’agilité, la testabilité et la capacité de déploiement sont impactées.

Architecture en pipeline

Architecture composée de blocs séquentiels où chaque étape traite et transmet les données à la suivante.

Points positifs :

-   Faible coût et simplicité, chaque bloc ayant un périmètre très limité
-   Architecture modulaire, testable, facile à faire évoluer et à composer

Points négatifs :

-   L’élasticité et la scalabilité sont faibles si l’architecture est déployée en monolithe, mais peuvent s’améliorer avec des unités de calcul éphémères (ex : Lambda)
-   La fiabilité est impactée lorsqu’un bloc échoue à s’exécuter

Architecture micro-noyau (Microkernel Architecture)

Architecture centrée sur un noyau minimal étendu par des plug-ins indépendants, idéale pour des applications personnalisables comme les IDE ou Jira.

Points positifs :

-   Faible coût et simplicité
-   La testabilité est relativement facile si un banc de test est mis en place pour le noyau principal, et que chaque plug-in reste relativement indépendant
-   Bonne déployabilité et fiabilité grâce à l’isolation des plug-ins

Points négatifs :

-   Scalabilité et tolérance aux pannes limitées

Architecture orientée domaine (Domain-Driven Architecture)

Hybride entre microservices et monolithe. Architecture découpée en domaines métiers indépendants mais déployés ensemble, combinant modularité et simplicité de gestion.

Il s’agit d’une architecture partitionnée par domaine, ce qui signifie que la structure est guidée par le métier plutôt que par des considérations purement techniques.

Elle combine une partie de la modularité des microservices avec la simplicité de gestion du monolithe.
C’est une approche naturelle dans un contexte de domain-driven design.

Les services sont des portions de l’application à gros grain, généralement appelées domaines (on parle ici de 4 à 12 services, contre des centaines dans une architecture microservices).
L’application est souvent déployée comme un monolithe (mas pas forcément). Idéalement, chaque domaine possède sa propre base de données. Il est aussi possible qu'une db puisse être consommé par plusieurs domaines.

Les services communiquent entre eux via le réseau en utilisant HTTP, gRPC ou des bus de messages (HTTP et gRPC peuvent être combinés dans un même service).

Points positifs :

-   Bien plus simple qu’une architecture microservices ou événementielle
-   Meilleure agilité, testabilité et déployabilité qu’un monolithe, ce qui accélère le time-to-market
-   Tolérance aux pannes et disponibilité élevées, car chaque service est déployé indépendamment

Points négatifs :

-   Scalabilité et élasticité bonnes (et meilleures qu’un monolithe), mais inférieures à celles d’une architecture microservices ; la simplicité, le trafic réseau et les coûts suivent la même logique
-   Une réflexion plus poussée est nécessaire par rapport à un monolithe, notamment pour les services devant être synchronisés.
    Le pattern Saga est largement utilisé pour répondre à ce défi.

Architecture événementielle (Event-Driven Architecture)

Système décentralisé réagissant à des événements, où les composants communiquent de manière asynchrone via des messages.

Les événements sont des réactions à des situations particulières.
Les deux topologies les plus connues sont le médiateur et le courtier (broker).

Le modèle broker offre un degré de réactivité plus élevé et un contrôle dynamique sur chaque événement.
Chaque unité de traitement est responsable du traitement du message et publie son résultat.
Ensuite, d’autres unités de traitement écoutent ces messages et les traitent si nécessaire.

Le modèle médiateur permet un meilleur contrôle du flux de traitement.
Il agit comme le gestionnaire d’état d’une requête.

Points positifs :

-   Performance, tolérance aux pannes et scalabilité sont les principaux avantages
-   L’ajout de nouvelles fonctionnalités est relativement simple

Points négatifs :

-   Un système équivalent est bien plus complexe à comprendre par rapport à d’autres topologies

Architecture microservices

Architecture distribuée où chaque service est autonome, faiblement couplé, et responsable d’un domaine fonctionnel précis.

Un des principaux défis est de trouver le bon niveau de granularité pour chaque service.
L’objectif est que chaque service reflète un processus métier unique.

Points positifs :

-   C’est l’architecture la plus scalable, élastique et évolutive

Points négatifs :

-   Complexité accrue pour comprendre l’ensemble du système et optimiser les performances
-   Nécessite de gérer la communication entre services et leur découverte
