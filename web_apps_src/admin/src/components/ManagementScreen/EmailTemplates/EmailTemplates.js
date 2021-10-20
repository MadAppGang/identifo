import React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { StaticFilesGeneralSection } from '~/components/shared/StaticFilesGeneralSection/StaticFilesGeneralSection';
import { getEmailTemplatesSettigns } from '~/modules/settings/selectors';
import { updateServerSettings } from '~/modules/settings/actions';


const EmailTemplates = () => {
  const dispatch = useDispatch();
  const settings = useSelector(getEmailTemplatesSettigns);

  return (
    <StaticFilesGeneralSection
      settings={settings}
      title="Email Templates"
      subtitle="These settings allow to specify paths to your email templates."
      onSubmit={values => dispatch(updateServerSettings({ emailTemplaits: values }))}
    />
  );
};

export default EmailTemplates;
