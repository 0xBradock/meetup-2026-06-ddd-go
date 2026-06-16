# Netscape Goes Bonkers

**Author:** Joel Spolsky
**Published:** November 20, 2000
**Source:** https://www.joelonsoftware.com/2000/11/20/netscape-goes-bonkers/

---

Netscape 6.0 serves as a cautionary tale about software development mistakes. The company made a critical strategic error by deciding to completely rewrite their browser from scratch — a decision that cost them three years of competitive ground against Internet Explorer. Lou Montulli, one of the original Navigator developers, confirmed this was "one of the major reasons I resigned from Netscape."

When Netscape finally released version 6.0, the product suffered from inadequate testing and numerous bugs. Rather than improving the codebase, the rewrite actually made things worse. As Montulli noted, old code "doesn't rust, it gets better" through bug fixes over time. The FTP code he had tuned over three years, though "ugly," worked reliably — something the rewritten version couldn't match.

The most significant complaint involved Netscape's decision to rebuild every UI widget from scratch to achieve cross-platform compatibility. This approach violated fundamental usability principles. The browser didn't support standard Windows conventions like right-click context menus, Alt+D keyboard shortcuts, Shift+Click for new windows, or scrolling mice. Most troubling for accessibility, the custom interface ignored system color schemes and font settings, creating barriers for vision-impaired users requiring high-contrast displays.

Without a prioritized development schedule, programmers implemented fun features rather than essential ones, leaving important functionality like standards compliance unfinished.

However, Spolsky acknowledged one strength: Netscape's Sidebar improved upon Microsoft Outlook's confusing design through effective use of real-world metaphors, making the interface more intuitive.

---

**Key quote:** Lou Montulli — *"It's one of the major reasons I resigned from Netscape. [...] They shouldn't have rewritten from scratch. They should have done this all in steps. Big chunky steps, fine, but steps."*
