package main

const myPRsQuery = `
query {
  viewer {
    login
    pullRequests(first: 30, states: OPEN, orderBy: {field: UPDATED_AT, direction: DESC}) {
      nodes {
        number
        title
        url
        isDraft
        mergeable
        reviewDecision
        baseRepository {
          nameWithOwner
        }
        author {
          login
        }
        createdAt
        updatedAt
        additions
        deletions
        commits(last: 1) {
          nodes {
            commit {
              statusCheckRollup {
                state
              }
            }
          }
        }
      }
    }
  }
}
`

const reviewRequestsQuery = `
query {
  search(query: "is:pr is:open review-requested:@me", type: ISSUE, first: 30) {
    nodes {
      ... on PullRequest {
        number
        title
        url
        isDraft
        author {
          login
        }
        baseRepository {
          nameWithOwner
        }
        createdAt
        updatedAt
        additions
        deletions
      }
    }
  }
}
`
