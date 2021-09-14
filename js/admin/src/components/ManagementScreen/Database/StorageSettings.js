import React, { useState, useEffect } from 'react';
import { useDispatch } from 'react-redux';
import EditIcon from '~/components/icons/EditIcon';
import LoadingIcon from '~/components/icons/LoadingIcon';
import Button from '~/components/shared/Button';
import SectionHeader from '~/components/shared/SectionHeader';
import { resetConnectionState, resetError } from '~/modules/database/actions';
import Form from './Form';
import Preview from './Preview';


const StorageSettings = (props) => {
  const {
    title, description, settings, activeTabIndex, onChange,
    progress, connectionState, postSettings, verifySettings } = props;
  const dispatch = useDispatch();
  const [editing, setEditing] = useState(false);

  const handleEditCancel = () => {
    dispatch(resetError());
    dispatch(resetConnectionState());
    setEditing(false);
  };

  const handlePostClick = async (data) => {
    postSettings(data);
  };

  const handleVerifyClick = (data) => {
    verifySettings(data);
  };

  useEffect(() => {
    return () => {
      dispatch(resetConnectionState());
    };
  }, []);

  return (
    <div className="iap-settings-section">
      <SectionHeader
        title={title}
        description={description}
      />
      <main>
        {editing && (
          <Form
            key={activeTabIndex}
            posting={!!progress}
            connectionStatus={connectionState}
            onChange={onChange}
            settings={settings}
            onSubmit={handlePostClick}
            onCancel={handleEditCancel}
            onVerify={handleVerifyClick}
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
