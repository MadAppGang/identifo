const createAppleService = ({ httpClient }) => {
    const uploadDevDomainAssociationFile = (file) => {
        const data = new FormData();
        const url = `${httpClient.getApiUrl()}/static/uploads/apple-domain-association`;

        data.append('file', file);

        return httpClient.post(url, data);
    };

    const uploadAppSiteAssociationFileContents = (fileContent) => {
        const url = `${httpClient.getApiUrl()}/static/template?name=apple-app-site-association`;
        const body = {
            contents: fileContent,
        };

        return httpClient.post(url, body);
    };

    const fetchAppSiteAssociationFileContents = async () => {
        const url = `${window.location.origin}/.well-known/apple-app-site-association`;
        const { data: blob } = await httpClient.get(url);
        const fileReader = new FileReader();

        return new Promise((resolve) => {
            fileReader.onload = () => {
                resolve(fileReader.result);
            };

            fileReader.readAsText(blob);
        });
    };

    return {
        fetchAppSiteAssociationFileContents,
        uploadAppSiteAssociationFileContents,
        uploadDevDomainAssociationFile,
    };
};

export default createAppleService;
