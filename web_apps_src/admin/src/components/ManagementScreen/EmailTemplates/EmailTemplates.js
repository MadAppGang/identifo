import React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { StaticFilesGeneralSection } from '~/components/shared/StaticFilesGeneralSection/StaticFilesGeneralSection';
import { getEmailTemplatesSettigns } from '~/modules/settings/selectors';
import { updateServerSettings } from '~/modules/settings/actions';
import { useVerification } from '~/hooks/useVerification';


const EmailTemplates = () => {
  const dispatch = useDispatch();
  const settings = useSelector(getEmailTemplatesSettigns);
  const [, verify] = useVerification();

  return (
    <StaticFilesGeneralSection
      settings={settings}
      title="Email Templates"
      subtitle="These settings allow to specify paths to your email templates."
      onSubmit={values => dispatch(updateServerSettings({ emailTemplaits: values }))}
      onVerify={values => dispatch(verify({ type: 'email_template_file_storage', ...values }))}
    />
  );
};

export default EmailTemplates;
