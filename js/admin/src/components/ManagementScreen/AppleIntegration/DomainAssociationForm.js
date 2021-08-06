import React, { useState } from 'react';
import FileInput from '~/components/shared/FileInput';
import Field from '~/components/shared/Field';
import Button from '~/components/shared/Button';
import SaveIcon from '~/components/icons/SaveIcon';
import LoadingIcon from '~/components/icons/LoadingIcon';
import useServices from '~/hooks/useServices';
import useProgressBar from '~/hooks/useProgressBar';
import useNotifications from '~/hooks/useNotifications';

const AppSiteAssociationForm = () => {
  const [file, setFile] = useState(null);
  const { progress, setProgress } = useProgressBar();
  const { notifySuccess, notifyFailure } = useNotifications();
  const services = useServices();

  const handleSubmit = async (event) => {
    event.preventDefault();

    if (!file) return;

    setProgress(70);

    try {
      await services.apple.uploadDevDomainAssociationFile(file);

      notifySuccess({
        title: 'Success',
        text: 'File has been uploaded.',
      });
    } catch (_) {
      notifyFailure({
        title: 'Something went wrong',
        text: 'File could not be uploaded.',
      });
    } finally {
      setProgress(100);
    }
  };

  return (
    <form className="domain-association-form" onSubmit={handleSubmit}>
      <Field
        label="Developer Domain Association File"
        subtext={file ? file.name : 'No file selected'}
      >
        <FileInput
          disablePathInput
          path="/.well-known/apple-developer-domain-association.txt"
          onFile={setFile}
        />
      </Field>

      <footer className="iap-apps-form__footer">
        <Button
          type="submit"
          Icon={progress ? LoadingIcon : SaveIcon}
          disabled={!file}
        >
          Upload
        </Button>
      </footer>
    </form>
  );
};

export default AppSiteAssociationForm;
