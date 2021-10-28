import React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { StaticFilesGeneralSection } from '~/components/shared/StaticFilesGeneralSection/StaticFilesGeneralSection';
import { updateServerSettings } from '~/modules/settings/actions';
import { getLoginWebAppSettings } from '~/modules/settings/selectors';
import { useVerification } from '~/hooks/useVerification';

const LoginWebAppSection = () => {
  const dispatch = useDispatch();
  const settings = useSelector(getLoginWebAppSettings);
  const [, verify] = useVerification();

  return (
    <StaticFilesGeneralSection
      settings={settings}
      title="Login Web App"
      subtitle="These settings allow to specify paths to your login web application."
      onSubmit={values => dispatch(updateServerSettings({ loginWebApp: values }))}
      onVerify={values => dispatch(verify({ type: 'spa_file_storage', ...values }))}
    />
  );
};

export default LoginWebAppSection;
