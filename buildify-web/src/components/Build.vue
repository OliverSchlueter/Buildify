<script lang="ts">
export default {
  props: [
    "baseLink",
    "isLatest",
    "id",
    "time",
    "hash",
    "message",
    "downloads",
    "fileName",
    "downloadLink"
  ],
  computed: {
    fullDownloadLink(){
        return this.baseLink + this.downloadLink;
    },
    shortHash() {
        return this.hash.substring(0, 7)
    },
    timeFormatted(){
        const t = new Date(this.time);
        return t.getDate() + "." + (t.getMonth() + 1) + "." + t.getFullYear() + " " + t.getHours() + ":" + t.getMinutes();
    }
  },
  methods: {
    download(){
        fetch("http://" + this.fullDownloadLink)
        .then(response => response.blob())
        .then(blob => {
            const url = URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = this.fileName;

            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);
            URL.revokeObjectURL(url);
        })
        .catch(error => {
            console.error('Error:', error);
    
            navigator.clipboard.writeText(this.fullDownloadLink);
            window.alert("Link to the file has been copied to your clipboard.\n" + this.fullDownloadLink)
        });
    }
  }
}
</script>

<template>
    <div :class="[isLatest ? 'isLatest' : '', 'build']">
        <div class="descriptions">
            <p class="id" @click="download">Build #{{ id }}</p>
            <div class="description">
                <p>
                    <span class="hash">{{ shortHash }}</span>
                    <span class="message">{{ message }}</span>
                </p>
            </div>
            <div class="description right">
                <p>
                    <span class="downloads">{{ downloads }} downloads</span>
                    <span class="time">{{ timeFormatted }}</span>
                </p>
            </div>
        </div>
    </div>
</template>

<style scoped>
    .build{
        padding: 0 20px;
        border-radius: 10px;
        display: flex;
        min-width: 600px;
        max-width: 800px;
        margin-left: auto;
        margin-right: auto;
    }

    .build:hover{
        background-color: rgba(8, 59, 66, 0.3);
    }

    .isLatest{
        margin-bottom: 20px;
    }

    .isLatest .id{
        background-color: rgb(10, 66, 8);
    }    

    .id{
        background-color: var(--primary-color-dark);
        padding: 5px 15px;
        border-radius: 15px;
        margin-right: 10px;
        font-weight: 700;
        cursor: pointer;
    }

    .id a{
        color: inherit;
        text-decoration: none;
    }

    .descriptions{
        display: flex;
        width: 100%;
        /* justify-content: space-between; */
    }

    .description{
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-content: center;
    }
    
    .descriptions .right{
        margin-left: auto;
    }
    .hash{
        color: var(--primary-color);
        margin-right: 10px;
        font-family: 'Courier New', Courier, monospace;
    }

    .downloads{
        color: rgb(175, 175, 175);
        margin-right: 30px;
    }
</style>