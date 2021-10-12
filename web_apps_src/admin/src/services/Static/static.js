const createStaticService = ({ httpClient }) => {
    const fetchStaticFile = async (name, ext = 'html') => {
        const url = `${httpClient.getApiUrl()}/static/template?name=${name}&ext=${ext}`;
        const response = await httpClient.get(url);

        return response.data.contents;
    };

    const updateStaticFile = async (name, ext, contents) => {
        const url = `${httpClient.getApiUrl()}/static/template?name=${name}&ext=${ext}`;
        await httpClient.put(url, { contents });
    };

    return {
        fetchStaticFile,
        updateStaticFile,
    };
};

export default createStaticService;
