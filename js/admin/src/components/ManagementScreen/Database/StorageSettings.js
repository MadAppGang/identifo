import React, { useState } from 'react';
import { useDispatch } from 'react-redux';
import Button from '~/components/shared/Button';
import EditIcon from '~/components/icons/EditIcon';
import LoadingIcon from '~/components/icons/LoadingIcon';
import SectionHeader from '~/components/shared/SectionHeader';
import { resetError } from '~/modules/database/actions';
import Preview from './Preview';
import Form from './Form';

const StorageSettings = (props) => {
  const { title, description, settings, progress, postSettings } = props;
  const dispatch = useDispatch();
  const [editing, setEditing] = useState(false);

  const handleEditCancel = () => {
    dispatch(resetError());
    setEditing(false);
  };

  const handlePostClick = async (data) => {
    postSettings(data);
  };

  return (
    <div className="iap-settings-section">
      <SectionHeader
        title={title}
        description={description}
      />
      <main>
        {editing && (
          <Form
            posting={!!progress}
            settings={settings}
            onSubmit={handlePostClick}
            onCancel={handleEditCancel}
          />
        )}
        {!editing && (
          <>
            <Preview fetching={progress} settings={settings} />
            <Button
              disabled={progress}
              Icon={progress ? LoadingIcon : EditIcon}
              onClick={() => setEditing(true)}
            >
              {`Edit ${title}`}
            </Button>
          </>
        )}
      </main>
    </div>
  );
};

export default StorageSettings;
