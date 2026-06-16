# Things You Should Never Do, Part I

**Author:** Joel Spolsky
**Published:** April 6, 2000
**Source:** https://www.joelonsoftware.com/2000/04/06/things-you-should-never-do-part-i/

---

Netscape 6.0 was entering its first public beta after a three-year gap since version 4.0. During those years, the company's market share had declined significantly. Joel Spolsky argues that this delay resulted from a critical strategic error: the decision to rewrite the entire codebase from scratch.

This mistake wasn't unique to Netscape. Borland made similar errors attempting to create dBase for Windows and rewriting Quattro Pro from scratch. Microsoft nearly made the same choice with a failed project called Pyramid, but fortunately maintained their existing codebase, which allowed them to continue shipping products.

Spolsky explains that programmers naturally want to rebuild systems entirely rather than incrementally improve them. However, this impulse stems from a fundamental principle: "It's harder to read code than to write it." Developers view existing code as messy and believe starting fresh will yield superior results. This assessment is typically incorrect.

Old code possesses inherent value through accumulated usage and testing. It contains countless bug fixes addressing real-world problems — edge cases in low-memory conditions, compatibility issues with older systems, and obscure interactions discovered only through extensive deployment. Those seemingly unnecessary lines represent weeks of troubleshooting and represent valuable institutional knowledge.

When companies discard existing code, they sacrifice market advantage by gifting competitors two to three years of development time. The rewriting process leaves them unable to ship new features or respond to market demands, essentially forcing a business shutdown during development.

Spolsky identifies three types of code problems that seem to justify rewrites:

- **Architectural issues** can be addressed through careful refactoring without abandoning the codebase.
- **Performance problems** typically affect only specific components and can be optimized surgically.
- **Cosmetic concerns** like inconsistent naming conventions require only minor automated fixes.

Crucially, starting fresh provides no assurance of better results. The original development team likely isn't involved, so experience doesn't transfer. New mistakes will emerge alongside reimplementation of old ones.

The outdated philosophy of "build one to throw away" proves dangerous at commercial scale. While experimental refactoring makes sense for individual functions or classes, abandoning entire systems represents catastrophic mismanagement. With proper oversight and software industry experience, companies like Netscape could have avoided this strategic self-sabotage.
