# Contributing to Gonsul

1. [Licence](#licence)
1. [Submitting a bug](#submitting-a-bug)
1. [Requesting or implementing a feature](#requesting-or-implementing-a-feature)
1. [Testing and submitting your changes](#testing-and-submitting-your-changes)
   1. [Code Style](#code-style)
   1. [Committing your changes](#committing-your-changes)
   1. [Pull requests and branching](#pull-requests-and-branching)
   1. [Credit](#credit)

## Licence

Gonsul is licenced under [The MIT LICENCE](LICENCE.md) for all code.

## Submitting a Bug

Bugs can be submitted to the [Github issue page](https://github.com/miniclip/gonsul/issues).

Gonsul is not perfect software and will be buggy. When submitting a bug, be
careful to know the following:

- The Go, Git and Consul versions you are running
- The flow you were attempting to use

You may be asked for further information regarding:

- Your environment, including any non-specified version,
  details about your operating system, and so on;
- Your project, including its structure, and possibly to remove build
  artifacts to start from a fresh build
- What it is you are trying to do exactly; we may provide alternative
  means to do so.

If you can provide an example code base to reproduce the issue on, we will
generally be able to provide more help, and faster.

Some contributors and maintainers may be unpaid developers
working on the project in their own time with limited resources. We
ask for respect and understanding and will try to provide the same back.

## Requesting or implementing a feature

Before requesting or implementing a new feature, please do the following:

- Verify in existing [issues](https://github.com/miniclip/gonsul/issues) whether
  the feature might already be in the works, or
  has already been rejected.
- Make sure you're using the latest release (or even the latest code, if you're
  going for _bleeding edge_)

If this is done, open up a issue. Tell us what is the feature you want,
why you need it, and why you think it should be in Gonsul itself.

We may discuss details with you regarding the implementation, and its inclusion
within the project.

We try to have as many of Gonsul's features tested as possible. Everything
that a user can do and should be repeatable in any way should be tested.

Upon creating a pull request you may be asked to add tests to the new/updated
behaviour, to guarantee backwards compatible where it should exist.

## Testing and submitting your changes

While we're not too formal when it comes to pull requests to the project,
we do appreciate users taking the time to conform to the guidelines that
follow.

We do expect all pull requests submitted to come with tests before
they are merged. If you cannot figure out how to write your tests properly, ask
in the pull request for guidance.

### Code Style

- Do not introduce trailing whitespace
- Indentation is 1 tab, not spaces.
- Try not to introduce lines longer than 100 characters
- Write small functions whenever possible, and use descriptive names for
  functions and variables.
- Comment tricky or non-obvious decisions made to explain their rationale.

### Committing your changes

It helps if your commits are structured as follows:

- Fixing a bug is one commit.
- Adding a feature is one commit.
- Adding two features is two commits.
- Two unrelated changes is two commits (and likely two Pull requests)

If you fix a (buggy) commit, squash (`git rebase -i`) the changes as a fixup
commit into the original commit, unless the patch was following a
maintainer's code review. In such cases, it helps to have separate commits.

The reviewer may ask you to later squash the commits together to provide
a clean commit history before merging in the feature.

It's important to write a proper commit title and description. The commit title
should be no more than 50 characters; it is the first line of the commit text. The
second line of the commit text must be left blank. The third line and beyond is
the commit message. You should write a commit message. If you do, wrap all
lines at 72 characters. You should explain what the commit does, what
references you used, and any other information that helps understanding your
changes.

### Pull requests and branching

All fixes to Gonsul end up requiring a +1 from one or more of the project's
maintainers. When opening a pull request, explain what the patch is doing
and if it makes sense, why the proposed implementation was chosen.

Try to use well-defined commits (one feature per commit) so that reading
them and testing them is easier for reviewers and while bisecting the code
base for issues.

During the review process, you may be asked to correct or edit a few things
before a final rebase to merge things. Do send edits as individual commits
to allow for gradual and partial reviews to be done by reviewers. Once the +1s
are given, rebasing is appreciated but not mandatory.

Please work in feature branches, and do not commit to `master` in your fork.

Provide a clean branch without merge commits.

If you can, pick a descriptive title for your pull request. This help for
eventual automated changelog generation.

### Credit

`gonsul` has been improved by
[many contributors](https://github.com/miniclip/gonsul/graphs/contributors)!
