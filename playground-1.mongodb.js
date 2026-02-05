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





// for each positionAssignment in all events, get the member details.
// if positionAssignment is null, default to an empty array.
// then group back to the original event structure.
// sort by event date ascending.
db.events.aggregate([
  {
    $unwind: {
      path: '$positionAssignments',
      preserveNullAndEmptyArrays: true
    }
  },
  {
    $lookup: {
      from: 'members',
      localField: 'positionAssignments.memberId',
      foreignField: '_id',
      as: 'memberDetails'
    }
  },
  {
    $unwind: {
      path: '$memberDetails',
      preserveNullAndEmptyArrays: true
    }
  },
  {
    $group: {
      _id: '$_id',
      name: { $first: '$name' },
      description: { $first: '$description' },
      date: { $first: '$date' },
      startTime: { $first: '$startTime' },
      endTime: { $first: '$endTime' },
      createdAt: { $first: '$createdAt' },
      updatedAt: { $first: '$updatedAt' },
      template: { $first: '$template' },
      // reconstruct positionAssignments array with member details
      positionAssignments: {
        $push: {
          _id: '$positionAssignments._id',
          description: '$positionAssignments.description',
          positionName: '$positionAssignments.positionName',
          member: '$memberDetails'
        }
      }
    }
  },
  {
    $sort: { date: 1 }
  }
]);

