# The Myth of the Software Rewrite

**Author:** Erik Dietrich
**Published:** October 18, 2015
**Source:** https://daedtech.com/the-myth-of-the-software-rewrite/

---

## The Setup

"We can't go on like this. We need to rewrite this thing from scratch."

These words infuriate CIOs and terrify managers — yet developers on the ground often say them emphatically. CIOs view a standing codebase as a paid-off asset. Developers live the day-to-day reality of a system leaving parts on the road after every pothole.

## Why Rewrites Fail

A software group starts productive, then hurries to ship features, makes a mess, always vowing to clean it up later. They never have the time. Eventually features slow to a crawl, developers are miserable, and attrition sets in. The developers want a total rewrite: "This time we know so many things we didn't know when we started — this time we'll get it right."

But won't the same thing be true in 3 more years? What makes you think giving the same group the same marching orders won't result in the same kind of code?

> "Insanity: doing the same thing over and over again and expecting different results." — Einstein

## The Alternative

There are legitimate cases for phasing out old software: hardware that's no longer manufactured, code in a defunct language nobody knows. But you don't need to rewrite software simply because developers made a mess of it while hurrying to meet deadlines.

The road back from a mess exists:

- Use automated tooling to identify and start improving the most dangerous parts of the code.
- **Automated tests are your friend** — characterize the system's current behavior with lots of automated tests and then work on refactoring.
- Bring in coaches or developers experienced with legacy rescues.
- Shift the team's priorities and help the business understand it's time to pay the piper on accumulated technical debt. They'll face a short-term slowdown to go faster sustainably over the long term — the same cost as a rewrite, but with an actual game-changing outcome.

## Conclusion

The rewrite is tempting when everyone is at wits' end. It's like making peace with a car payment and getting excited about the leather seats of the luxury thing you're going to buy. But software isn't a car.

> **"The software is a mess because the group made it a mess, and it'll only get and stay clean if the group cleans it."**
