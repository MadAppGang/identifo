import React from 'react';
import { isPhone, getInitials } from '~/utils';
import PhoneIcon from '~/components/icons/PhoneIcon.svg';

const UserIcon = (props) => {
  const { username, email, phone } = props;

  let contents;

  if (!username && !email) {
    if (phone) {
      contents = <PhoneIcon className="user-phone-icon" />;
    } else {
      contents = '-';
    }
  } else if (isPhone(username)) {
    contents = <PhoneIcon className="user-phone-icon" />;
  } else {
    contents = getInitials(username, email);
  }

  return (
    <div className="iap-datagrid-row__user-icon">
      {contents}
    </div>
  );
};

export default UserIcon;
