# Welcome to PubDoc

## TL;DR
The goal of this repository is twofold:
1. To make documentation better for the reader.
2. To make _creating_ documentation a better experience for the writer.

## How to Use
This repository is currently very "bare bones." To run the tool and interface with Confluence,  
you may need to contact an IT administrator to obtain permissions to access Confluence data.  

Visit [this link](https://id.atlassian.com/manage-profile/security/api-tokens) to create your API token.  
Next, create a `.env` file with the following contents:

### `.env` File Example
```env
CONFLUENCE_USERNAME=<username>
CONFLUENCE_API_KEY=<api key>
CONFLUENCE_DOMAIN=<your company>.atlassian.net
CONFLUENCE_SPACEID=<A space you don't care about littering in>
```
To get a spaceID of a space. Go to the space then click on "space settings" in the top left. Then run this in the browser to get the spaceID. `https://YOURURL.atlassian.net/wiki/rest/api/space/KEY`

## Pontifications & Justifications
If you're interested in why we're working on this, here's some context:  

Want to feel depressed? Go look at a Confluence documentation page. You can see who has visited or read the page.  
You'll likely notice that the number of viewers per page is remarkably low. Who's to blame—the reader or the writer?  
Who's to say. Regardless, one thing is certain: the documentation could be better.  

By making our documentation _programmatic_, we can improve it in several ways:
1. We can ensure that documentation is "correct" and "clear."
    - Often, the same information is included in multiple locations. Using _programmatic_ documentation,  
    we can enforce consistency across all locations.
    - _Assuming we place this in the monorepo,_ the documentation will be **directly** tied to a Git commit.
2. Documentation could include a "_programmatic_" representation of behavior.  
See [this link](https://flexgen.atlassian.net/wiki/x/CIAzPw) for further details.
3. And much more.
## To-Do
- [x] Create a basic API interface for Confluence documentation.
- [ ] Define a _single source of truth_ (SSOT) for "variable"-based documentation.
    - [ ] Develop a mechanism for recurring access to the SSOT.
- [ ] Investigate adding Matplotlib or other graphing libraries to provide interactive "hover" HTML elements.
- [ ] Etc.
