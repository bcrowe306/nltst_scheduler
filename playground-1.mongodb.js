/* global use, db */
// MongoDB Playground
// To disable this template go to Settings | MongoDB | Use Default Template For Playground.
// Make sure you are connected to enable completions and to be able to run a playground.
// Use Ctrl+Space inside a snippet or a string literal to trigger completions.
// The result of the last command run in a playground is shown on the results panel.
// By default the first 20 documents will be returned with a cursor.
// Use 'console.log()' to print to the debug output.
// For more documentation on playgrounds please refer to
// https://www.mongodb.com/docs/mongodb-vscode/playgrounds/

// Select the database to use.
use('nltst_scheduler');





// Aggregation to get events grouped by each member. Should return an array of members with their associated events and their positionAssignments for each event.
// The memberId exist in the positionAssignments array inside each event.
// example output:
// [
//   {
//     _id: "449ad960-87ee-4636-a0bf-a43cba94db38",
//     firstName: "Brandon",
//     lastName: "Crowe",
//     email: "bcrowe306@gmail.com",
//     phoneNumber: "8136062719",
//     createdAt: {
//       $date: "2026-02-03T23:21:46.923Z"
//     },
//     updatedAt: {
//       $date: "2026-02-03T23:24:23.919Z"
//     },
//     events: [
//       {
//         _id: "08acfa6b-b839-4f04-bfea-ab6f739cf8ca",
//         name: "Sunday AM Service",
//         description: "Normal Sunday AM worship service",
//         template: "",
//         startTime: "11:00",
//         endTime: "13:00",
//         date: {
//           $date: "2026-02-08T00:00:00Z"
//         },
//         reminderInterval: 0,
//         reminderEnabled: false,
//         teamId: "76b445c2-f2bf-406a-80a5-203648fce254",
//         createdAt: {
//           $date: "2026-02-05T21:34:35.634Z"
//         },
//         updatedAt: {
//           $date: "2026-02-05T21:34:57.351Z"
//         },
//         positionName: "Usher", // from positionAssignments only that match the memberId
//       }
      
//     ]

//   }
// ]

var memberId = "449ad960-87ee-4636-a0bf-a43cba94db38"; // Example memberId to filter events.positionAssignments for a specific member
var startDate = new Date("2026-02-010T00:00:00Z"); // Example start date for filtering events
var endDate = new Date("2026-03-10T23:59:59Z"); // Example end date for filtering events

db.members.aggregate([
  // Match specific member by memberId
  {
    $match: { _id: memberId } 
  },
  {
    $lookup: {
      from: "events",
      let: { memberId: "$_id" },
      pipeline: [
        {
          $match: {
            $expr: {
              $in: [
                "$$memberId",
                {
                  $map: {
                    input: "$positionAssignments",
                    as: "pa",
                    in: "$$pa.memberId"
                  }
                }
              ]
            }
          },
        },

        // Filter events by date range
        {
          $match: {
            date: { $gte: startDate, $lte: endDate }
          }
        },
        {
          $addFields: {
            positionName: {
              $let: {
                vars: {
                  matchedPA: {
                    $arrayElemAt: [
                      {
                        $filter: {
                          input: "$positionAssignments",
                          as: "pa",
                          cond: { $eq: ["$$pa.memberId", "$$memberId"] }
                        }
                      },
                      0
                    ]
                  }
                },
                in: "$$matchedPA.positionName"
              }
            }
          }
        },
        {
          $project: {
            positionAssignments: 0
          }
        },
        { 
          $sort: { date: 1, startTime: 1 } // Sort events by date ascending
        },
      ],
      as: "events"
    },
  }
]);
