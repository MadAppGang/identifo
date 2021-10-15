import { pause } from '~/utils';

const createAppleServiceMock = () => {
  const uploadDevDomainAssociationFile = async () => {
    await pause(400);
  };

  const uploadAppSiteAssociationFileContents = async () => {
    await pause(400);
  };

  const fetchAppSiteAssociationFileContents = async () => {
    await pause(400);
    return '{\n\t\n\t\n}';
  };

  return {
    uploadDevDomainAssociationFile,
    uploadAppSiteAssociationFileContents,
    fetchAppSiteAssociationFileContents,
  };
};

export default createAppleServiceMock;
