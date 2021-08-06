const createApplicationService = ({ httpClient }) => {
    const fetchApplications = async () => {
        const url = `${httpClient.getApiUrl()}/apps`;
        const { data = [] } = await httpClient.get(url);

        return data;
    };

    const fetchApplicationById = async (id) => {
        const url = `${httpClient.getApiUrl()}/apps/${id}`;
        const { data } = await httpClient.get(url);

        return data;
    };

    const postApplication = async (application) => {
        const url = `${httpClient.getApiUrl()}/apps`;
        const { data } = await httpClient.post(url, application);

        return data;
    };

    const alterApplication = async (id, changes) => {
        const url = `${httpClient.getApiUrl()}/apps/${id}`;
        const { data } = await httpClient.put(url, changes);

        return data;
    };

    const deleteApplicationById = async (id) => {
        const url = `${httpClient.getApiUrl()}/apps/${id}`;
        const { data } = await httpClient.delete(url);

        return data;
    };

    const fetchFederatedLoginProviders = async () => {
        const url = `${httpClient.getApiUrl()}/federated-providers`;
        const { data } = await httpClient.get(url);
        return data;
    };

    return Object.freeze({
        fetchApplications,
        fetchApplicationById,
        postApplication,
        alterApplication,
        deleteApplicationById,
        fetchFederatedLoginProviders,
    });
};

export default createApplicationService;
