import { curry } from './fn';

export const months = [
  'Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul","Aug', 'Sep', 'Oct', 'Nov', 'Dec',
];

export const sameDay = (timestampA, timestampB) => {
  const dateA = new Date(timestampA);
  const dateB = new Date(timestampB);

  const sameDate = dateA.getDate() === dateB.getDate();
  const sameMonth = dateA.getMonth() === dateB.getMonth();
  const sameYear = dateA.getFullYear() === dateB.getFullYear();

  return sameDate && sameMonth && sameYear;
};

export const isToday = curry(sameDay, new Date().getTime());

export const formatDateForTable = (timestamp) => {
  if (!timestamp) {
    return 'Never';
  }

  if (isToday(timestamp)) {
    return 'Today';
  }

  const date = new Date(timestamp);

  return `${months[date.getMonth()]} ${date.getDate()}`;
};
