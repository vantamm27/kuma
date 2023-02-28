<template>
    <div class="mb-3">
        <label for="lark-url" class="form-label">{{ $t("Post URL") }}</label>
        <input id="lark-url" v-model="$parent.notification.larkURL" type="url" pattern="https?://.+" class="form-control" required>
    </div>

    <div class="mb-3">
        <label for="lark-content-type" class="form-label">{{ $t("Content Type") }}</label>
        <select id="lark-content-type" v-model="$parent.notification.larkContentType" class="form-select" required>
            <option value="json">
                application/json
            </option>
            <option value="form-data">
                multipart/form-data
            </option>
        </select>

        <div class="form-text">
            <p>{{ $t("larkJsonDesc", ["\"application/json\""]) }}</p>
            <i18n-t tag="p" keypath="larkFormDataDesc">
                <template #multipart>"multipart/form-data"</template>
                <template #decodeFunction>
                    <strong>json_decode($_POST['data'])</strong>
                </template>
            </i18n-t>
        </div>
    </div>
</template>
